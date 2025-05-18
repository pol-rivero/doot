class Doot < Formula
  desc     "Fast, simple and intuitive dotfiles manager that just gets the job done"
  homepage "https://github.com/pol-rivero/doot"
  version  "{{VERSION}}"
  license  "MIT"
  head     "https://github.com/pol-rivero/doot.git", branch: "main"

  depends_on "git"
  depends_on "git-crypt"

  on_macos do
    on_arm do
      url "https://github.com/pol-rivero/doot/releases/download/{{VERSION}}/doot-darwin-arm64"
      sha256 "{{DARWIN_ARM64_CHECKSUM}}"
    end
    on_intel do
      url "https://github.com/pol-rivero/doot/releases/download/{{VERSION}}/doot-darwin-x86_64"
      sha256 "{{DARWIN_X86_CHECKSUM}}"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/pol-rivero/doot/releases/download/{{VERSION}}/doot-linux-arm64"
      sha256 "{{LINUX_ARM64_CHECKSUM}}"
    end
    on_intel do
      url "https://github.com/pol-rivero/doot/releases/download/{{VERSION}}/doot-linux-x86_64"
      sha256 "{{LINUX_X86_CHECKSUM}}"
    end
  end

  def install
    mv Dir["doot-*"].first, "doot"
    chmod 0755, "doot"
    bin.install "doot"

    generate_completions_from_executable(bin/"doot", "completion")
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/doot --version")
  end
end
