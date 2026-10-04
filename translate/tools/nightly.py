#!/usr/bin/env python3
"""Resolve a published Wave nightly or install a previously locked snapshot.

Bootstrap tooling only; the translation service and HTTP implementation are Wave.
Use a persistent --store to retain archives after upstream replaces the nightly.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import tarfile
import tempfile
import urllib.request

API = 'https://api.github.com/repos/wavefnd/Wave'
DOWNLOAD = 'https://github.com/wavefnd/Wave/releases/download/nightly/'
DEFAULT_LOCK = Path(__file__).resolve().parents[1] / 'wave-toolchain.lock.json'


def open_url(url, accept='application/vnd.github+json'):
    headers = {'Accept': accept, 'User-Agent': 'Wave-Translation-bootstrap'}
    token = os.environ.get('GH_TOKEN')
    if token and url.startswith('https://api.github.com/'):
        headers['Authorization'] = 'Bearer ' + token
    return urllib.request.urlopen(urllib.request.Request(url, headers=headers), timeout=60)


def fetch(url, accept='application/vnd.github+json'):
    with open_url(url, accept) as response:
        return response.read()


def valid_name(name):
    return isinstance(name, str) and re.fullmatch(r'[A-Za-z0-9._-]+', name) is not None


def verify(data, size, digest):
    if len(data) != size or hashlib.sha256(data).hexdigest() != digest:
        raise ValueError('Artifact size or SHA-256 mismatch')


def resolve(target):
    before = json.loads(fetch(API + '/releases/tags/nightly'))
    manifests = [a for a in before['assets'] if a['name'].startswith('nightly-') and a['name'].endswith('.json') and a['state'] == 'uploaded']
    if len(manifests) != 1:
        raise ValueError('Nightly publication is incomplete or ambiguous; retry after publishing completes')
    asset = manifests[0]
    raw = fetch(API + '/releases/assets/' + str(asset['id']), 'application/octet-stream')
    digest = asset.get('digest', '')
    if not re.fullmatch(r'sha256:[0-9a-f]{64}', digest):
        raise ValueError('Nightly manifest has no published SHA-256')
    verify(raw, asset['size'], digest[7:])
    manifest = json.loads(raw)
    if manifest['schema_version'] != 1 or not re.fullmatch('[0-9a-f]{40}', manifest['source_sha']):
        raise ValueError('Unsupported nightly manifest')
    generation = manifest['generation']
    if not valid_name(generation) or not generation.startswith(manifest['source_sha'] + '-'):
        raise ValueError('Invalid nightly generation')
    if asset['name'] != 'nightly-' + generation + '.json':
        raise ValueError('Manifest filename does not match its generation')
    name = 'wave-nightly-' + generation + '-' + target + '.tar.gz'
    record = manifest['assets'][name]
    archives = [a for a in before['assets'] if a['name'] == name and a['state'] == 'uploaded']
    if len(archives) != 1 or archives[0]['size'] != record['size'] or archives[0].get('digest') != 'sha256:' + record['sha256']:
        raise ValueError('Archive metadata does not match the manifest')
    after = json.loads(fetch(API + '/releases/tags/nightly'))
    identity = lambda r: (r['id'], r['target_commitish'], sorted((a['id'], a['name'], a.get('digest'), a['size']) for a in r['assets']))
    if identity(before) != identity(after) or before['target_commitish'] != manifest['source_sha']:
        raise ValueError('Nightly changed during resolution; retry')
    return {
        'schema_version': 1, 'source_sha': manifest['source_sha'],
        'generation': generation, 'compiler_version': manifest['compiler_version'],
        'target': target,
        'manifest': {'name': asset['name'], 'sha256': digest[7:], 'size': asset['size']},
        'archive': {'name': name, 'sha256': record['sha256'], 'size': record['size']},
    }, raw


def read_lock(path):
    lock = json.loads(path.read_text())
    if lock['schema_version'] != 1 or not re.fullmatch('[0-9a-f]{40}', lock['source_sha']):
        raise ValueError('Unsupported lock file')
    generation = lock['generation']
    if not valid_name(generation) or not generation.startswith(lock['source_sha'] + '-'):
        raise ValueError('Invalid generation')
    if lock['target'] not in ('x86_64-linux-gnu', 'aarch64-linux-gnu'):
        raise ValueError('This foundation currently validates Linux targets only')
    if not isinstance(lock['compiler_version'], str) or not lock['compiler_version']:
        raise ValueError('Invalid compiler version')
    for key in ('manifest', 'archive'):
        item = lock[key]
        if not valid_name(item['name']) or not re.fullmatch('[0-9a-f]{64}', item['sha256']) or not isinstance(item['size'], int) or item['size'] <= 0:
            raise ValueError('Invalid locked artifact')
    if lock['manifest']['name'] != 'nightly-' + generation + '.json' or lock['archive']['name'] != 'wave-nightly-' + generation + '-' + lock['target'] + '.tar.gz':
        raise ValueError('Locked names do not match the generation')
    return lock


def verify_file(path, item):
    digest = hashlib.sha256()
    size = 0
    with path.open('rb') as source:
        while chunk := source.read(1024 * 1024):
            size += len(chunk)
            digest.update(chunk)
    if size != item['size'] or digest.hexdigest() != item['sha256']:
        raise ValueError('Artifact size or SHA-256 mismatch: ' + item['name'])


def retain(store, item, supplied=None, offline=False):
    destination = store / item['sha256'] / item['name']
    if destination.exists():
        verify_file(destination, item)
        return destination
    if offline and supplied is None:
        raise ValueError('Snapshot is not retained locally: ' + str(destination))
    destination.parent.mkdir(parents=True, exist_ok=True)
    with tempfile.NamedTemporaryFile(dir=destination.parent, delete=False) as temp:
        temporary = Path(temp.name)
        try:
            if supplied is not None:
                temp.write(supplied)
            else:
                with open_url(DOWNLOAD + item['name'], 'application/octet-stream') as response:
                    total = 0
                    while chunk := response.read(1024 * 1024):
                        total += len(chunk)
                        if total > item['size']:
                            raise ValueError('Download exceeds locked size')
                        temp.write(chunk)
            temp.flush()
            verify_file(temporary, item)
            temporary.replace(destination)
        finally:
            temporary.unlink(missing_ok=True)
    return destination


def validate_manifest(lock, path):
    manifest = json.loads(path.read_text())
    if manifest['schema_version'] != 1 or any(manifest[key] != lock[key] for key in ('source_sha', 'generation', 'compiler_version')):
        raise ValueError('Retained manifest does not match the lock')
    record = manifest['assets'][lock['archive']['name']]
    if any(record[key] != lock['archive'][key] for key in ('sha256', 'size')):
        raise ValueError('Retained manifest does not match the archive')


def extract(archive, destination, lock):
    if destination.exists():
        raise ValueError('Install destination already exists; use a new directory')
    destination.parent.mkdir(parents=True, exist_ok=True)
    staging = Path(tempfile.mkdtemp(prefix='.wave-stage-', dir=destination.parent))
    try:
        with tarfile.open(archive, 'r:gz') as tar:
            # data_filter rejects escaping paths and links and unsafe file types.
            if not hasattr(tarfile, 'data_filter'):
                raise ValueError('Python with tarfile.data_filter is required (3.12+ recommended)')
            tar.extractall(staging, filter='data')
        compilers = list(staging.glob('*/wavec'))
        if len(compilers) != 1 or not compilers[0].is_file():
            raise ValueError('Archive must contain one Wave package')
        package = compilers[0].parent
        if not (package / 'llvm').is_dir() or not (package / 'std').is_dir():
            raise ValueError('Wave package is missing its bundled LLVM or standard library')
        (package / 'wave-translation-install.json').write_text(json.dumps(lock, indent=2) + '\n')
        package.rename(destination)
    finally:
        shutil.rmtree(staging)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest='command', required=True)
    update = sub.add_parser('resolve', help='Write a candidate lock and retain its manifest and archive')
    update.add_argument('--target', choices=['x86_64-linux-gnu', 'aarch64-linux-gnu'], default='x86_64-linux-gnu')
    update.add_argument('--lock', type=Path, required=True, help='Candidate output; never overwrites an existing lock')
    update.add_argument('--store', type=Path, required=True, help='Persistent archive store chosen by the operator')
    install = sub.add_parser('install', help='Install exactly the locked snapshot; never resolves latest')
    install.add_argument('--lock', type=Path, default=DEFAULT_LOCK)
    install.add_argument('--store', type=Path, required=True)
    install.add_argument('--destination', type=Path, required=True)
    install.add_argument('--offline', action='store_true')
    args = parser.parse_args()
    if args.command == 'resolve':
        if args.lock.exists():
            parser.error('Candidate lock already exists; choose a new path')
        lock, raw = resolve(args.target)
        retain(args.store, lock['manifest'], raw)
        retain(args.store, lock['archive'])
        args.lock.parent.mkdir(parents=True, exist_ok=True)
        with args.lock.open('x') as output:
            output.write(json.dumps(lock, indent=2) + '\n')
        print(args.lock)
    else:
        lock = read_lock(args.lock)
        manifest = retain(args.store, lock['manifest'], offline=args.offline)
        validate_manifest(lock, manifest)
        archive = retain(args.store, lock['archive'], offline=args.offline)
        extract(archive, args.destination, lock)
        print(args.destination / 'wavec')


if __name__ == '__main__':
    try:
        main()
    except (ValueError, KeyError, OSError, tarfile.TarError) as error:
        raise SystemExit(str(error))
