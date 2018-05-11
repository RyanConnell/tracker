http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.12.0/rules_go-0.12.0.tar.gz",
    sha256 = "c1f52b8789218bb1542ed362c4f7de7052abcf254d865d96fb7ba6d44bc15ee3",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains",)
go_rules_dependencies()
go_register_toolchains()
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

# External Dependencies

load("@bazel_gazelle//:deps.bzl", "go_repository")

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    tag = "v1.6.1",
)

go_repository(
    name = "com_github_gorilla_context",
    importpath = "github.com/gorilla/context",
    tag = "v1.1",
)

go_repository(
    name = "com_github_gorilla_sessions",
    importpath = "github.com/gorilla/sessions",
    tag = "v1.1",
)

go_repository(
    name = "com_github_gorilla_securecookie",
    importpath = "github.com/gorilla/securecookie",
    tag = "v1.1.1",
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    remote = "git@github.com:golang/oauth2",
    vcs = "git",
    commit = "cdc340f7c179dbbfa4afd43b7614e8fcadde4269",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    remote = "git@github.com:golang/net",
    vcs = "git",
    commit = "f73e4c9ed3b7ebdd5f699a16a880c2b1994e50dd",
)

go_repository(
    name = "com_google_cloud_go",
    importpath = "cloud.google.com/go",
    tag = "v0.22.0",
)

go_repository(
    name = "com_github_go-sql-driver_mysql",
    importpath = "github.com/go-sql-driver/mysql",
    tag = "v1.3",
)