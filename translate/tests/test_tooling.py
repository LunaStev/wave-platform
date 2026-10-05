#!/usr/bin/env python3
"""Offline regression tests for nightly replacement and isolated builds."""
import copy
import hashlib
import io
import json
from pathlib import Path
import sys
import tarfile
import tempfile
import unittest
from unittest.mock import patch

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / 'tools'))
import build
import nightly


def artifact(name, data):
    return {'name': name, 'size': len(data), 'sha256': hashlib.sha256(data).hexdigest()}


def snapshot():
    sha = 'a' * 40
    generation = sha + '-123-1'
    archive = artifact('wave-nightly-' + generation + '-x86_64-linux-gnu.tar.gz', b'archive')
    manifest = {
        'schema_version': 1, 'source_sha': sha, 'generation': generation,
        'compiler_version': 'test', 'assets': {archive['name']: archive},
    }
    raw = json.dumps(manifest).encode()
    lock = {key: manifest[key] for key in ('schema_version', 'source_sha', 'generation', 'compiler_version')}
    lock.update(target='x86_64-linux-gnu', archive=archive,
                manifest=artifact('nightly-' + generation + '.json', raw))
    release = {'id': 1, 'target_commitish': sha, 'assets': [
        {'id': index, 'name': item['name'], 'size': item['size'],
         'digest': 'sha256:' + item['sha256'], 'state': 'uploaded'}
        for index, item in enumerate([lock['manifest'], archive], start=1)
    ]}
    return lock, raw, release


class ToolingTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.lock, self.raw, self.release = snapshot()
        self.lock_path = self.root / 'lock.json'
        self.lock_path.write_text(json.dumps(self.lock))

    def resolve(self, before=None, after=None, raw=None):
        before = before if before is not None else self.release
        after = after if after is not None else before
        with patch.object(nightly, 'fetch', side_effect=[json.dumps(before), raw or self.raw, json.dumps(after)]):
            return nightly.resolve('x86_64-linux-gnu')

    def test_resolves_one_coherent_generation(self):
        self.assertEqual(self.resolve(), (self.lock, self.raw))
        self.assertEqual(nightly.read_lock(self.lock_path), self.lock)

    def test_replaced_release_or_assets_are_rejected(self):
        for change in ['release', 'asset', 'source']:
            changed = copy.deepcopy(self.release)
            if change == 'release':
                changed['id'] += 1
            elif change == 'asset':
                changed['assets'][1]['id'] += 1
            else:
                changed['target_commitish'] = 'b' * 40
            with self.subTest(change=change), self.assertRaisesRegex(ValueError, 'changed'):
                self.resolve(after=changed)

    def test_partial_publication_and_wrong_digests_fail(self):
        cases = []
        missing = copy.deepcopy(self.release)
        missing['assets'].pop()
        cases.append(missing)
        mismatch = copy.deepcopy(self.release)
        mismatch['assets'][1]['digest'] = 'sha256:' + '0' * 64
        cases.append(mismatch)
        unpublished = copy.deepcopy(self.release)
        unpublished['assets'][0]['state'] = 'new'
        cases.append(unpublished)
        for release in cases:
            with self.subTest(release=release), self.assertRaises(ValueError):
                self.resolve(before=release)
        with self.assertRaisesRegex(ValueError, 'mismatch'):
            self.resolve(raw=b'tampered')

    def test_published_source_must_match_manifest(self):
        changed = copy.deepcopy(self.release)
        changed['target_commitish'] = 'b' * 40
        with self.assertRaisesRegex(ValueError, 'changed'):
            self.resolve(before=changed)

    def test_lock_rejects_paths_and_cross_generation_names(self):
        for field, value in [('name', '../escape'), ('name', 'different.tar.gz'),
                             ('sha256', 'bad'), ('size', 0)]:
            changed = copy.deepcopy(self.lock)
            changed['archive'][field] = value
            self.lock_path.write_text(json.dumps(changed))
            with self.subTest(field=field, value=value), self.assertRaises(ValueError):
                nightly.read_lock(self.lock_path)

    def test_retained_snapshot_works_after_upstream_disappears(self):
        store = self.root / 'store'
        nightly.retain(store, self.lock['manifest'], self.raw)
        nightly.retain(store, self.lock['archive'], b'archive')
        with patch.object(nightly, 'open_url', side_effect=AssertionError('No network allowed')):
            manifest = nightly.retain(store, self.lock['manifest'], offline=True)
            nightly.validate_manifest(self.lock, manifest)
            self.assertEqual(nightly.retain(store, self.lock['archive'], offline=True).read_bytes(), b'archive')

    def test_missing_or_corrupted_store_never_falls_forward(self):
        item = self.lock['archive']
        with patch.object(nightly, 'open_url', side_effect=AssertionError('No network allowed')):
            with self.assertRaisesRegex(ValueError, 'not retained'):
                nightly.retain(self.root, item, offline=True)
            path = nightly.retain(self.root, item, b'archive')
            path.write_bytes(b'corrupt')
            with self.assertRaisesRegex(ValueError, 'mismatch'):
                nightly.retain(self.root, item)

    def test_failed_download_leaves_no_published_or_temporary_file(self):
        for data in [b'wrong!!', b'oversized response']:
            with patch.object(nightly, 'open_url', return_value=io.BytesIO(data)), self.assertRaises(ValueError):
                nightly.retain(self.root, self.lock['archive'])
            self.assertEqual(list((self.root / self.lock['archive']['sha256']).iterdir()), [])

    def test_streamed_download_is_verified_and_retained(self):
        with patch.object(nightly, 'open_url', return_value=io.BytesIO(b'archive')):
            self.assertEqual(nightly.retain(self.root, self.lock['archive']).read_bytes(), b'archive')

    def test_install_checks_manifest_lock_binding(self):
        path = nightly.retain(self.root, self.lock['manifest'], self.raw)
        wrong = copy.deepcopy(self.lock)
        wrong['compiler_version'] = 'other'
        with self.assertRaisesRegex(ValueError, 'does not match'):
            nightly.validate_manifest(wrong, path)
        wrong = copy.deepcopy(self.lock)
        wrong['archive']['sha256'] = '0' * 64
        with self.assertRaisesRegex(ValueError, 'does not match'):
            nightly.validate_manifest(wrong, path)

    def archive(self, include_std=True, unsafe=False):
        path = self.root / 'package.tar.gz'
        with tarfile.open(path, 'w:gz') as tar:
            names = ['package/wavec', 'package/llvm/lib/placeholder']
            if include_std:
                names.append('package/std/placeholder')
            if unsafe:
                names.append('../escaped')
            for name in names:
                member = tarfile.TarInfo(name)
                member.size = 1
                tar.addfile(member, io.BytesIO(b'x'))
        return path

    def test_extract_requires_complete_package_and_preserves_existing_install(self):
        destination = self.root / 'installed'
        with self.assertRaisesRegex(ValueError, 'standard library'):
            nightly.extract(self.archive(include_std=False), destination, self.lock)
        self.assertFalse(destination.exists())
        nightly.extract(self.archive(), destination, self.lock)
        self.assertEqual(json.loads((destination / 'wave-translation-install.json').read_text()), self.lock)
        with self.assertRaisesRegex(ValueError, 'already exists'):
            nightly.extract(self.archive(), destination, self.lock)

    def test_escaping_archive_path_is_rejected(self):
        with self.assertRaises(tarfile.FilterError):
            nightly.extract(self.archive(unsafe=True), self.root / 'installed', self.lock)
        self.assertFalse((self.root / 'escaped').exists())
        self.assertFalse((self.root / 'installed').exists())

    @patch.object(build.platform, 'system', return_value='Linux')
    @patch.object(build.platform, 'machine', return_value='x86_64')
    def test_build_uses_locked_std_and_rejects_other_install(self, *_):
        toolchain = self.root / 'toolchain'
        nightly.extract(self.archive(), toolchain, self.lock)
        with patch.object(build.subprocess, 'run') as run:
            build.build(toolchain, self.lock_path, self.root / 'server', self.root / 'build')
            command = run.call_args.args[0]
            self.assertEqual(command[:3], [str(toolchain / 'wavec'), '--std-root', str(toolchain / 'std')])
            receipt = toolchain / 'wave-translation-install.json'
            receipt.write_text('{}')
            with self.assertRaisesRegex(ValueError, 'differs'):
                build.build(toolchain, self.lock_path, self.root / 'server', self.root / 'build')
            self.assertEqual(run.call_count, 1)


if __name__ == '__main__':
    unittest.main()
