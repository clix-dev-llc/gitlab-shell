.PHONY: validate verify verify_ruby verify_golang test test_ruby test_golang setup _install build compile check clean

validate: verify test

verify: verify_ruby verify_golang

verify_ruby:
	bundle exec rubocop

verify_golang:
	support/go-format check

test: test_ruby test_golang

test_ruby:
	# bin/gitlab-shell must exist and needs to be the Ruby version for
	# rspec to be able to test.
	cp bin/gitlab-shell-ruby bin/gitlab-shell
	# bin/gitlab-shell-authorized-keys-check and bin/gitlab-shell-authorized-principals-check
	# should link to ruby scripts for rspec to be able to test.
	ln -sf ./gitlab-shell-authorized-keys-check-ruby bin/gitlab-shell-authorized-keys-check
	ln -sf ./gitlab-shell-authorized-principals-check-ruby bin/gitlab-shell-authorized-principals-check
	bundle exec rspec --color --tag '~go' --format d spec
	rm -f bin/gitlab-shell
	ln -sf ./gitlab-shell bin/gitlab-shell-authorized-keys-check
	ln -sf ./gitlab-shell bin/gitlab-shell-authorized-principals-check

test_golang:
	support/go-test

setup: _install bin/gitlab-shell

_install:
	bin/install

build: bin/gitlab-shell
compile: bin/gitlab-shell
bin/gitlab-shell:
	bin/compile

check:
	bin/check

clean:
	rm -f bin/gitlab-shell