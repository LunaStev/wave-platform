package blog

import "testing"

func TestContentLanguageIgnoresSourceCode(t *testing.T) {
	for _, test := range []struct{ content, language string }{
		{"The Wave compiler translates your source code into native machine code. This release improves the standard library and fixes compiler errors.", "en"},
		{"Wave 프로그래밍 언어의 새로운 소식을 전합니다. 이번 릴리즈에는 컴파일러 개선과 표준 라이브러리 기능이 포함되어 있습니다. 개발에 참여해 주신 모든 분께 감사드립니다.", "ko"},
		{"Wave プログラミング言語の新しいリリースを紹介します。この更新ではコンパイラの改善と標準ライブラリの新しい機能が追加されました。開発に参加してくださった皆様に感謝します。", "ja"},
		{"v0.2.2", "en"},
	} {
		post := Post{Content: test.content + "\n```wave\nvar result: i32 = compile(source);\n```"}
		if got := post.ContentLanguage(); got != test.language {
			t.Errorf("got %s, want %s for %q", got, test.language, test.content)
		}
	}
}
