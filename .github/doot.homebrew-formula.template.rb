class Doot < Formula
  desc     "A fast, simple and intuitive dotfiles manager that just gets the job done"
  homepage "https://github.com/pol-rivero/doot"
  license  "MIT"
  version  "{{VERSION}}"
  head     "https://github.com/pol-rivero/doot.git", branch: "main"
  url      "https://github.com/pol-rivero/doot/archive/refs/tags/{{VERSION}}.tar.gz"
  sha256   "{{TARBALL_CHECKSUM}}"

  depends_on "git"
  depends_on "git-crypt"
  depends_on "go" => :build


  on_macos do
    on_arm do
      resource "binary" do
        url "https://github.com/pol-rivero/doot/releases/download/{{VERSION}}/doot-darwin-arm64"
        sha256 "{{DARWIN_ARM64_CHECKSUM}}"
      end
    end
    on_intel do
      resource "binary" do
        url "https://github.com/pol-rivero/doot/releases/download/{{VERSION}}/doot-darwin-x86_64"
        sha256 "{{DARWIN_X86_CHECKSUM}}"
      end
    end
  end

  on_linux do
    on_arm do
      resource "binary" do
        url "https://github.com/pol-rivero/doot/releases/download/{{VERSION}}/doot-linux-arm64"
        sha256 "{{LINUX_ARM64_CHECKSUM}}"
      end
    end
    on_intel do
      resource "binary" do
        url "https://github.com/pol-rivero/doot/releases/download/{{VERSION}}/doot-linux-x86_64"
        sha256 "{{LINUX_X86_CHECKSUM}}"
      end
    end
  end

  resource "source" do
    url "https://github.com/pol-rivero/doot/archive/refs/tags/{{VERSION}}.tar.gz"
    sha256 "{{TARBALL_CHECKSUM}}"
  end


  def install
    if build.from_source?
      resource("source").stage do
        system "go", "build", "-o", "doot", "-trimpath", "-mod=readonly", "-modcacherw"
        bin.install "doot"
      end
    else
      resource("binary").stage do
        # rename the downloaded binary to "doot" and install
        bin.install Dir["doot*"] => "doot"
      end
    end

    generate_completions_from_executable(bin/"doot", "completion")
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/doot --version")
  end
end
