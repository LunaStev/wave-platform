import assert from 'node:assert/strict'
import test from 'node:test'
import { parseToolchainCatalog, formatToolchainSize, toolchainChecksumCommand } from '../src/services/toolchains.ts'
const filename = 'wave-llvm-21.1.8-linux-riscv64-r1.tar.xz'
const valid = { target: 'linux-riscv64', llvm_version: '21.1.8', revision: 'r1', filename,
  size_bytes: 1048576, sha256: 'a'.repeat(64), published_at: '2026-10-03T00:00:00Z',
  url: `https://wave-lang.dev/downloads/toolchains/llvm/21.1.8/r1/${filename}` }
const catalog = (entry: unknown) => ({ schema_version: 1, bundles: [entry] })
test('only verified catalog URLs become same-origin downloads', () => {
  assert.equal(parseToolchainCatalog(catalog(valid))[0].url, `/downloads/toolchains/llvm/21.1.8/r1/${filename}`)
  assert.equal(formatToolchainSize(valid.size_bytes), '1.0 MiB')
  assert.deepEqual(parseToolchainCatalog({ schema_version: 1, bundles: [] }), [])
})
test('checksum instructions follow each architecture, version and revision', () => {
  const loongFilename = 'wave-llvm-21.1.9-linux-loong64-r2.tar.xz'
  const bundles = parseToolchainCatalog({ schema_version: 1, bundles: [valid, {
    ...valid, target: 'linux-loong64', llvm_version: '21.1.9', revision: 'r2', filename: loongFilename,
    url: `https://wave-lang.dev/downloads/toolchains/llvm/21.1.9/r2/${loongFilename}`,
  }] })
  assert.deepEqual(bundles.map(toolchainChecksumCommand), [
    'sha256sum --check wave-llvm-21.1.8-linux-riscv64-r1.tar.xz.sha256',
    'sha256sum --check wave-llvm-21.1.9-linux-loong64-r2.tar.xz.sha256',
  ])
  assert.equal(bundles[1].url, `/downloads/toolchains/llvm/21.1.9/r2/${loongFilename}`)
})
test('invalid hashes, paths, duplicate files, versions and sizes are rejected', () => {
  for (const change of [{ url: 'javascript:alert(1)' }, { url: 'https://outside.invalid/file' },
    { filename: '../secret' }, { sha256: 'pending' }, { size_bytes: 0 }, { size_bytes: '100' },
    { target: 'linux-amd64' }, { llvm_version: '22.1.0' }, { revision: '../r1' }, { published_at: 'unknown' }]) {
    assert.throws(() => parseToolchainCatalog(catalog({ ...valid, ...change })))
  }
  assert.throws(() => parseToolchainCatalog({ schema_version: 1, bundles: [valid, valid] }))
})
