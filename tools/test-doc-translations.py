#!/usr/bin/env python3
"""Regression checks for incremental translation contributions."""
import contextlib
import importlib.util
import io
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('translations', Path(__file__).with_name('check-doc-translations.py'))
checker = importlib.util.module_from_spec(spec)
spec.loader.exec_module(checker)


class TranslationChecks(unittest.TestCase):
    def test_coverage_policy(self):
        self.assertEqual(checker.coverage_errors({'a.md'}, {}, 'pt'), [])
        self.assertEqual(checker.coverage_errors({'a.md'}, {}, 'en'), ['a.md'])
        self.assertEqual(checker.coverage_errors({'a.md'}, {'wrong.md'}, 'pt'), ['wrong.md'])

    def run_fixture(self, target, change_code=False):
        with tempfile.TemporaryDirectory() as directory:
            docs = Path(directory)
            source = '''---
translation_set_id: sample
path: language/sample
locale: ko
group: language
group_order: 2
order: 1
title: Sample
summary: Sample
---
## Example
```wave
var value: i32 = 1;
```
[Next](/docs/ko/language/next)
'''
            translated = source.replace('locale: ko', 'locale: pt').replace('/docs/ko/language/next', target)
            if change_code:
                translated = translated.replace('i32 = 1', 'i32 = 2')
            for locale, text in [('ko', source), ('pt', translated)]:
                (docs / locale / 'language').mkdir(parents=True)
                (docs / locale / 'language/sample.md').write_text(text)
            (docs / 'en/language').mkdir(parents=True)
            (docs / 'en/language/next.md').write_text('English fallback')
            with patch.object(checker, 'DOCS', docs), patch('sys.argv', ['check', '--locale', 'pt']), contextlib.redirect_stdout(io.StringIO()):
                checker.main()

    def test_link_to_untranslated_page_uses_existing_english_target(self):
        self.run_fixture('/docs/pt/language/next')

    def test_missing_fallback_target_is_rejected(self):
        with self.assertRaises(SystemExit):
            self.run_fixture('/docs/pt/language/missing')

    def test_link_must_keep_selected_language(self):
        with self.assertRaises(SystemExit):
            self.run_fixture('/docs/en/language/next')

    def test_partial_translation_still_validates_code(self):
        with self.assertRaises(SystemExit):
            self.run_fixture('/docs/pt/language/next', change_code=True)


if __name__ == '__main__':
    unittest.main()
