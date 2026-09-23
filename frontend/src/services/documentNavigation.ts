export type DocumentationProject = 'wave' | 'whale'

export interface NavigationDocument {
  path: string
  groupOrder: number
  order: number
}

export function documentationProject(path: string): DocumentationProject {
  return path === 'whale' || path.startsWith('whale/') ? 'whale' : 'wave'
}

export function documentationCatalog(locale: string, project: DocumentationProject): string {
  return `/docs/${locale}${project === 'whale' ? '/whale' : ''}`
}

export function documentationPath(path: string): string {
  return path === 'whale' ? '' : path
}

export function projectDocuments<T extends NavigationDocument>(translated: T[], english: T[], project: DocumentationProject): T[] {
  const byPath = new Map<string, T>()
  for (const item of [...english, ...translated]) {
    if ((item.path.startsWith('whale/') ? 'whale' : 'wave') === project) byPath.set(item.path, item)
  }
  return [...byPath.values()].sort((a, b) => a.groupOrder - b.groupOrder || a.order - b.order || (a.path < b.path ? -1 : a.path > b.path ? 1 : 0))
}
