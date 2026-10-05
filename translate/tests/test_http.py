#!/usr/bin/env python3
"""Exercise the real Wave server using split TCP writes and malformed requests."""
import json
import os
from pathlib import Path
import socket
import struct
import subprocess
import time
import unittest

BINARY = Path(os.environ.get('WAVE_TRANSLATE_BINARY', '/tmp/wave-translation-server'))


class HTTPFoundation(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        with socket.socket() as available:
            available.bind(('127.0.0.1', 0))
            cls.port = available.getsockname()[1]
        cls.process = subprocess.Popen([str(BINARY)], env={**os.environ, 'WAVE_TRANSLATE_PORT': str(cls.port)}, stdout=subprocess.DEVNULL, stderr=subprocess.PIPE)
        for _ in range(100):
            if cls.process.poll() is not None:
                raise RuntimeError(cls.process.stderr.read().decode())
            try:
                with socket.create_connection(('127.0.0.1', cls.port), timeout=.1):
                    return
            except OSError:
                time.sleep(.02)
        cls.process.terminate()
        cls.process.wait(timeout=5)
        raise RuntimeError('Wave server did not listen')

    @classmethod
    def tearDownClass(cls):
        cls.process.terminate()
        cls.process.wait(timeout=5)
        cls.process.stderr.close()

    def exchange(self, fragments, delay=0):
        with socket.create_connection(('127.0.0.1', self.port), timeout=4) as connection:
            for fragment in fragments:
                connection.sendall(fragment)
                if delay:
                    time.sleep(delay)
            data = b''
            while True:
                try:
                    part = connection.recv(4096)
                except ConnectionResetError:
                    break
                if not part:
                    break
                data += part
        head, body = data.split(b'\r\n\r\n', 1)
        lines = head.split(b'\r\n')
        headers = dict(line.split(b': ', 1) for line in lines[1:])
        self.assertEqual(headers[b'Connection'], b'close')
        return int(lines[0].split()[1]), headers, body

    def request(self, path='/healthz', method='GET', headers='Host: localhost\r\n'):
        return self.exchange([f'{method} {path} HTTP/1.1\r\n{headers}\r\n'.encode()])

    def test_health_and_truthful_readiness(self):
        status, headers, body = self.request()
        self.assertEqual(status, 200)
        self.assertEqual(int(headers[b'Content-Length']), len(body))
        self.assertEqual(json.loads(body)['service'], 'wave-translation')
        status, _, body = self.request('/readyz')
        self.assertEqual(status, 503)
        self.assertEqual(json.loads(body), {'ready': False, 'reason': 'model_not_configured'})

    def test_head_retains_content_length_without_body(self):
        for path in ['/healthz', '/readyz', '/missing']:
            expected = self.request(path)
            status, headers, body = self.request(path, 'HEAD')
            self.assertEqual(status, expected[0])
            self.assertEqual(headers[b'Content-Length'], expected[1][b'Content-Length'])
            self.assertEqual(body, b'')

    def test_tcp_fragmentation_and_case_insensitive_headers(self):
        pieces = [b'GE', b'T /healthz HTTP/1.1\r', b'\nhOsT: localhost\r\nContent-Length: 0\r\n\r', b'\n']
        self.assertEqual(self.exchange(pieces, .02)[0], 200)

    def test_head_errors_have_no_body(self):
        for headers in ['', 'Host: a\r\nContent-Length: 1\r\n',
                        'Host: a\r\nTransfer-Encoding: chunked\r\n']:
            expected = self.request(headers=headers)
            status, response_headers, body = self.request(method='HEAD', headers=headers)
            self.assertEqual(status, expected[0])
            self.assertEqual(response_headers[b'Content-Length'], expected[1][b'Content-Length'])
            self.assertEqual(body, b'')

    def test_invalid_framing(self):
        for headers in ['', 'Host: localhost\r\nHost: other\r\n', 'Host : localhost\r\n',
                        'Host: \r\n', 'Host: a,b\r\n', 'Host: localhost\r\n Folded: no\r\n',
                        'Host: localhost\r\nContent-Length: 0\r\nContent-Length: 0\r\n',
                        'Host: localhost\r\nContent-Length: -1\r\n',
                        'Host: localhost\r\nContent-Length: 0\r\nTransfer-Encoding: chunked\r\n',
                        'Host: localhost\r\nBad: x\x00y\r\n']:
            with self.subTest(headers=headers):
                self.assertEqual(self.request(headers=headers)[0], 400)
        self.assertEqual(self.request(headers='Host: localhost\r\nTransfer-Encoding: chunked\r\n')[0], 501)
        self.assertEqual(self.request(headers='Host: localhost\r\nContent-Length: 999999999999999999999999\r\n')[0], 413)

    def test_invalid_request_lines(self):
        for line in [b' GET /healthz HTTP/1.1', b'GET  /healthz HTTP/1.1', b'GET /bad\x00path HTTP/1.1', b'GET /bad#fragment HTTP/1.1']:
            self.assertEqual(self.exchange([line + b'\r\nHost: localhost\r\n\r\n'])[0], 400)
        self.assertEqual(self.exchange([b'GET /healthz HTTP/1.0\r\nHost: localhost\r\n\r\n'])[0], 505)

    def test_limits_and_routes(self):
        self.assertEqual(self.request('/missing')[0], 404)
        status, headers, _ = self.request(method='POST')
        self.assertEqual(status, 405)
        self.assertEqual(headers[b'Allow'], b'GET, HEAD')
        self.assertEqual(self.exchange([b'GET /healthz HTTP/1.1\r\nHost: a\r\nX: ' + b'a' * 8192])[0], 431)
        self.assertEqual(self.request(headers='Host: a\r\n' + 'X: a\r\n' * 100)[0], 431)
        # One connection serves one request, including when two arrive together.
        status, _, body = self.exchange([b'GET /healthz HTTP/1.1\r\nHost: a\r\n\r\nGET /readyz HTTP/1.1\r\nHost: a\r\n\r\n'])
        self.assertEqual(status, 200)
        self.assertNotIn(b'HTTP/1.1', body)

    def test_exact_header_limit_and_head_overflow(self):
        prefix = b'GET /healthz HTTP/1.1\r\nHost: a\r\nX: '
        request = prefix + b'a' * (8192 - len(prefix) - 4) + b'\r\n\r\n'
        self.assertEqual(self.exchange([request])[0], 200)
        self.assertEqual(self.exchange([request[:-4] + b'a\r\n\r\n'])[0], 431)
        status, _, body = self.exchange([b'HEAD /healthz HTTP/1.1\r\nHost: a\r\nX: ' + b'a' * 8192])
        self.assertEqual(status, 431)
        self.assertEqual(body, b'')

    def test_client_disconnect_does_not_stop_server(self):
        with socket.create_connection(('127.0.0.1', self.port), timeout=4) as connection:
            connection.sendall(b'GET /healthz HTTP/1.1\r\nHost: a\r\n\r\n')
            connection.setsockopt(socket.SOL_SOCKET, socket.SO_LINGER, struct.pack('ii', 1, 0))
        self.assertEqual(self.request()[0], 200)

    def test_slow_request_has_total_deadline(self):
        started = time.monotonic()
        self.assertEqual(self.exchange([b'GET ', b'/healthz '], .6)[0], 408)
        self.assertLess(time.monotonic() - started, 3.5)
        self.assertEqual(self.request()[0], 200)

    def test_invalid_port_fails_before_listening(self):
        for port in ['0', '65536', '-1', 'abc']:
            result = subprocess.run([str(BINARY)], env={**os.environ, 'WAVE_TRANSLATE_PORT': port}, capture_output=True, timeout=3)
            self.assertEqual(result.returncode, 2)


if __name__ == '__main__':
    unittest.main()
