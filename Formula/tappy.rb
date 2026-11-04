class Tappy < Formula
  desc "Terminal clicking-speed game"
  homepage "https://github.com/plyr4/tappy"
  head "https://github.com/plyr4/tappy.git", branch: "main"

  depends_on "go" => :build
  depends_on :macos

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w"), "./cmd/tappy"
  end

  test do
    assert_path_exists bin/"tappy"
  end
end
