#!/usr/bin/env python3
"""Check translation coverage, Markdown structure, examples, and navigation.

Korean is the canonical source. These checks detect mechanical translation
regressions; technical meaning and natural phrasing still need editorial review.
"""
from pathlib import Path, PurePosixPath
import re

ROOT = Path(__file__).resolve().parents[1]
DOCS = ROOT / 'wavedoc'
LOCALES = ['en', 'ja', 'zh', 'es', 'de', 'ru', 'id', 'vi']
FENCE = re.compile(r'^```[^\n]*\n.*?^```[ \t]*$', re.M | re.S)
MARKER = re.compile(r'<!-- wave-example: ([a-z0-9-]+) -->')
LINK = re.compile(r'\]\(([^\s)]+)\)')
SHARED_METADATA = ['translation_set_id', 'path', 'group', 'group_order', 'order']
UNICODE_LITERALS = ['`한`', '`"한"`', "`'한'`"]


def metadata(text):
    return dict(line.split(': ', 1) for line in text.split('---', 2)[1].strip().splitlines())


def structure(text):
    prose = FENCE.sub('', text.split('---', 2)[2])
    return {
        'headings': re.findall(r'^(#{2,6}) ', prose, re.M),
        'tables': [len(re.split(r'(?<!\\)\|', line))
                   for line in prose.splitlines() if line.startswith('|')],
        'list_items': len(re.findall(r'^\s*(?:[-*]|\d+\.) ', prose, re.M)),
    }


def main():
    source = {str(path.relative_to(DOCS / 'ko')): path.read_text()
              for path in (DOCS / 'ko').rglob('*.md')}
    if not source:
        raise SystemExit('No Korean source documents found')
    failures = []
    for locale in LOCALES:
        actual = {str(path.relative_to(DOCS / locale)): path
                  for path in (DOCS / locale).rglob('*.md')}
        if set(actual) != set(source):
            failures.append((locale, 'coverage', sorted(set(source) ^ set(actual))))
        for name, korean in source.items():
            if name not in actual:
                continue
            translated = actual[name].read_text()
            original_meta, translated_meta = metadata(korean), metadata(translated)
            for key in SHARED_METADATA:
                if original_meta[key] != translated_meta.get(key):
                    failures.append((locale, name, 'metadata ' + key))
            if translated_meta.get('locale') != locale:
                failures.append((locale, name, 'locale metadata'))
            if not translated_meta.get('title') or not translated_meta.get('summary'):
                failures.append((locale, name, 'missing title or summary'))
            if FENCE.findall(korean) != FENCE.findall(translated):
                failures.append((locale, name, 'code/command/output changed'))
            if MARKER.findall(korean) != MARKER.findall(translated):
                failures.append((locale, name, 'example ID mismatch'))
            if structure(korean) != structure(translated):
                failures.append((locale, name, 'headings/table/list structure changed'))

            prose = FENCE.sub('', translated)
            without_literals = prose
            for literal in UNICODE_LITERALS:
                without_literals = without_literals.replace(literal, '')
            if re.search('[가-힣]', without_literals):
                failures.append((locale, name, 'untranslated Korean'))
            if re.search(r'⟦[UTSC]\d+⟧', prose):
                failures.append((locale, name, 'unrestored placeholder'))
            for target in LINK.findall(prose):
                target = target.split('#', 1)[0]
                if not target or '://' in target or target.startswith('mailto:'):
                    continue
                if target.startswith('/docs/'):
                    parts = target.removeprefix('/docs/').split('/', 1)
                    if parts[0] != locale:
                        failures.append((locale, name, 'link leaves selected locale: ' + target))
                        continue
                    path = parts[1] if len(parts) > 1 else ''
                    if path not in ('', 'whale', 'stdlib') and path + '.md' not in actual:
                        failures.append((locale, name, 'broken link: ' + target))
                elif not target.startswith('/'):
                    path = str(PurePosixPath(name).parent / target)
                    if path + '.md' not in actual:
                        failures.append((locale, name, 'broken relative link: ' + target))
        print(f'{locale}: {len(actual)} documents')

    if failures:
        for issue in failures:
            print('FAIL', *issue)
        raise SystemExit(f'{len(failures)} translation checks failed')
    print(f'PASS: {len(source)} source documents, {len(source) * len(LOCALES)} translations; '
          'code, structure and links match.')


if __name__ == '__main__':
    main()
