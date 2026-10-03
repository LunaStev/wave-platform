export interface ToolchainBundle {
  target: 'linux-riscv64' | 'linux-loong64'
  llvm_version: string
  revision: string
  filename: string
  size_bytes: number
  sha256: string
  url: string
  published_at: string
}

export function parseToolchainCatalog(value: unknown): ToolchainBundle[] {
  if (!value || typeof value !== 'object') throw new Error('Invalid toolchain catalog')
  const catalog = value as Record<string, unknown>
  if (catalog.schema_version !== 1 || !Array.isArray(catalog.bundles)) throw new Error('Invalid toolchain catalog')
  const seen = new Set<string>()
  return catalog.bundles.map((raw: unknown) => {
    if (!raw || typeof raw !== 'object') throw new Error('Invalid toolchain entry')
    const b = raw as ToolchainBundle
    if (!['linux-riscv64', 'linux-loong64'].includes(b.target)
        || typeof b.llvm_version !== 'string' || !/^21\.\d+\.\d+$/.test(b.llvm_version)
        || typeof b.revision !== 'string' || !/^r[1-9]\d*$/.test(b.revision)
        || typeof b.sha256 !== 'string' || !/^[a-f0-9]{64}$/.test(b.sha256)
        || !Number.isSafeInteger(b.size_bytes) || b.size_bytes <= 0
        || typeof b.published_at !== 'string' || !Number.isFinite(Date.parse(b.published_at))) {
      throw new Error('Invalid toolchain metadata')
    }
    const filename = `wave-llvm-${b.llvm_version}-${b.target}-${b.revision}.tar.xz`
    const path = `/downloads/toolchains/llvm/${b.llvm_version}/${b.revision}/${filename}`
    if (b.filename !== filename || b.url !== `https://wave-lang.dev${path}` || seen.has(path)) {
      throw new Error('Invalid or duplicate toolchain download')
    }
    seen.add(path)
    // Serve through this platform's origin, including local previews.
    return { ...b, url: path }
  })
}

export function formatToolchainSize(bytes: number): string {
  return `${(bytes / 1048576).toFixed(1)} MiB`
}

export function toolchainChecksumCommand(bundle: ToolchainBundle): string {
  return `sha256sum --check ${bundle.filename}.sha256`
}
