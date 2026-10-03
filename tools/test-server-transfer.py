#!/usr/bin/env python3
"""Exercise transfer scripts with temporary installations and a fake Docker CLI."""
import os
from pathlib import Path
import shutil
import subprocess
import tarfile
import tempfile
import unittest

REPOSITORY = Path(__file__).resolve().parents[1]


class ServerTransferTests(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="wave-transfer-test-")
        self.addCleanup(self.temporary.cleanup)
        self.root = Path(self.temporary.name)
        binaries = self.root / "bin"
        binaries.mkdir()
        docker = binaries / "docker"
        docker.write_text('#!/bin/sh\nprintf "%s\\n" "$*" >> "$TRANSFER_TEST_LOG"\n')
        docker.chmod(0o755)
        self.log = self.root / "docker.log"
        self.environment = {**os.environ, "PATH": f"{binaries}:{os.environ['PATH']}",
                            "TRANSFER_TEST_LOG": str(self.log)}

    def installation(self, name):
        destination = self.root / name
        destination.mkdir()
        for script in ("export-server.sh", "import-server.sh"):
            shutil.copy2(REPOSITORY / script, destination / script)
        return destination

    def run_script(self, installation, script, archive):
        return subprocess.run(["bash", str(installation / script), str(archive)],
                              env=self.environment, capture_output=True, text=True)

    def export(self, with_toolchains):
        source = self.installation("source")
        (source / ".env").write_text("TRANSFER_TEST_ONLY=1\n")
        (source / "data").mkdir()
        (source / "data" / "record").write_text("private fixture")
        if with_toolchains:
            (source / "toolchains").mkdir()
            (source / "toolchains" / "index.json").write_text('{"schema_version":1,"bundles":[]}')
            (source / "toolchains" / "sdk.tar.xz").write_bytes(b"SDK fixture")
        archive = self.root / "transfer.tar.gz"
        result = self.run_script(source, "export-server.sh", archive)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(archive.stat().st_mode & 0o777, 0o600)
        return archive

    def test_round_trip_includes_public_downloads(self):
        archive = self.export(with_toolchains=True)
        with tarfile.open(archive) as contents:
            self.assertIn("toolchains/sdk.tar.xz", contents.getnames())
        destination = self.installation("destination")
        (destination / "toolchains").mkdir()
        (destination / "toolchains" / ".gitkeep").touch()
        result = self.run_script(destination, "import-server.sh", archive)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual((destination / "toolchains" / "sdk.tar.xz").read_bytes(), b"SDK fixture")
        self.assertEqual((destination / "data" / "record").read_text(), "private fixture")
        self.assertEqual((destination / ".env").stat().st_mode & 0o777, 0o600)

    def test_legacy_archive_without_downloads(self):
        archive = self.export(with_toolchains=False)
        destination = self.installation("destination")
        result = self.run_script(destination, "import-server.sh", archive)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertTrue((destination / "data" / "record").is_file())

    def test_fresh_clone_with_tracked_downloads(self):
        archive = self.export(with_toolchains=True)
        destination = self.installation("destination")
        (destination / "toolchains").mkdir()
        (destination / "toolchains" / "sdk.tar.xz").write_bytes(b"SDK fixture")
        subprocess.run(["git", "init", "-q", str(destination)], check=True)
        subprocess.run(["git", "-C", str(destination), "add", "toolchains"], check=True)
        subprocess.run(["git", "-C", str(destination), "-c", "user.name=Transfer Test",
                        "-c", "user.email=test@example.invalid", "commit", "-qm", "Fixture"], check=True)
        result = self.run_script(destination, "import-server.sh", archive)
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual((destination / "toolchains" / "sdk.tar.xz").read_bytes(), b"SDK fixture")

    def test_existing_downloads_are_not_overwritten(self):
        archive = self.export(with_toolchains=True)
        destination = self.installation("destination")
        (destination / "toolchains").mkdir()
        existing = destination / "toolchains" / "sdk.tar.xz"
        existing.write_bytes(b"existing download")
        previous_calls = self.log.read_text()
        result = self.run_script(destination, "import-server.sh", archive)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("existing toolchain downloads", result.stderr)
        self.assertEqual(existing.read_bytes(), b"existing download")
        self.assertFalse((destination / ".env").exists())
        self.assertEqual(self.log.read_text(), previous_calls)


if __name__ == "__main__":
    unittest.main()
