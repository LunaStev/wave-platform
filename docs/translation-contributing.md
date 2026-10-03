# Contributing a documentation translation

Korean (`wavedoc/ko/`) is the canonical source. English (`wavedoc/en/`) is a
reading aid and the published fallback. The website interface remains English
and Korean; these locales apply to documentation only.

## Language and writing system

Use one language code per language, without regional directories, URLs, or
translation variants. Use `pt`, never `pt-BR` or `pt-PT`, and apply the same rule
to other languages. Preserve normal accents and language-specific spelling;
Latin script does not mean ASCII or transliteration.

| Code | Language | Writing system |
| --- | --- | --- |
| en | English | Latin |
| ko | Korean | Korean |
| ja | Japanese | Normal Japanese writing: kanji, hiragana, katakana |
| zh | Chinese | Simplified Chinese |
| es | Spanish | Latin |
| de | German | Latin |
| ru | Russian | Cyrillic |
| id | Indonesian–Malay | Latin; one shared documentation locale |
| vi | Vietnamese | Latin |
| pt | Portuguese | Latin |
| fr | French | Latin |
| pl | Polish | Latin |
| nl | Dutch | Latin |
| tr | Turkish | Latin |
| it | Italian | Latin |

Use broadly understandable wording within each language, and keep terminology
consistent with existing translations. For a future language with multiple
writing systems, select one before opening translation tasks; for example,
Serbo-Croatian would use Latin. This example does not enable another locale.

`wavedoc/locales.json` records supported locales, writing systems, and coverage
requirements. Open Graph requires a territory in its metadata values; that does
not create a separate documentation locale or impose a regional translation.

## One document per contribution

- Choose an open translation issue and check that nobody is already working on it.
- Read its Korean source and the corresponding English page. Ask about unclear
  technical meaning in the issue before changing a language or API contract.
- Create `wavedoc/{locale}/{same-directory}/{same-file}.md`. Keep the source's
  `translation_set_id`, `path`, `group`, `group_order`, and `order`; set `locale`
  to the selected language code. Translate `title`, `summary`, and all prose.
- Preserve headings, tables, lists, exercises, answers, example IDs, and fenced
  code/commands/output. Do not translate API identifiers or change examples.
- Keep documentation links in the selected locale. A link may lead to an English
  fallback while its target is awaiting translation. Update translated heading
  fragments to match the destination page; check them in the rendered document.
- Run the checks below, then open a PR referencing the issue. Request a fluent
  review for terminology and meaning. Passing checks cannot certify prose quality.

For example, when contributing Portuguese:

```sh
python3 tools/check-doc-translations.py --locale pt
python3 tools/check-doc-examples.py --links-only
git diff --check
```

The six new locales (`pt`, `fr`, `pl`, `nl`, `tr`, `it`) permit partial coverage.
Missing pages use the English version. Existing complete locales retain their
coverage requirement. Every added translation is still checked for structure,
metadata, code, and links. A new locale may therefore pass with zero documents;
that means its foundation is valid, not that translation is complete.

Do not add copied source text, placeholders, or empty Markdown pages to make a
locale appear translated. The initial `.gitkeep` files only preserve directories
and are not published. Add completed documents gradually. Only actual published
translations receive localized detail entries in the sitemap and `hreflang`.

No per-document route registration is needed. Rebuild and restart the application
after changing Markdown because it is embedded in the Go binary. See
[document authoring](document-authoring.md) for preview and metadata details.
