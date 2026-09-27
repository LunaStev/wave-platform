import redirects from '../../../wavedoc/redirects.json' with { type: 'json' }

export type DocumentationProject = 'wave' | 'stdlib' | 'whale'

export interface NavigationDocument {
  path: string
  groupOrder: number
  order: number
}

export function documentationProject(path: string): DocumentationProject {
  if (path === 'whale' || path.startsWith('whale/')) return 'whale'
  if (path === 'stdlib' || path.startsWith('stdlib/')) return 'stdlib'
  const libraryPaths = ['reference/standard-library', 'reference/string-and-bytes', 'reference/memory-and-buffer', 'reference/system-io-network-process']
  return libraryPaths.includes(path) ? 'stdlib' : 'wave'
}

export function documentationCatalog(locale: string, project: DocumentationProject): string {
  return `/docs/${locale}${project === 'wave' ? '' : `/${project}`}`
}

export function documentationPath(path: string): string {
  return path === 'whale' || path === 'stdlib' ? '' : path
}

export function projectDocuments<T extends NavigationDocument>(translated: T[], english: T[], project: DocumentationProject): T[] {
  const byPath = new Map<string, T>()
  for (const item of [...english, ...translated]) {
    if (documentationProject(item.path) === project) byPath.set(item.path, item)
  }
  return [...byPath.values()].sort((a, b) => a.groupOrder - b.groupOrder || a.order - b.order || (a.path < b.path ? -1 : a.path > b.path ? 1 : 0))
}

export function canonicalDocumentPath(path: string): string {
  return (redirects as Record<string, string>)[path] ?? path
}
