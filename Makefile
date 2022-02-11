default: install

build_dir=build
dest_dir=$(build_dir)/$(shell sed -En "s/^baseurl: \"(.*)\"/\1/gp" _config.yaml)

.PHONY: help
help:
	@egrep '(^\S)|^$$' Makefile

.PHONY: install
install:
	bundle config set --local path vendor/bundle
	bundle install

.PHONY: upgrade
upgrade:
	bundle update

.PHONY: build
build:
	go run ../cmd/nmstatectl --help > user-guide/main-help.txt
	go run ../cmd/nmstatectl gen --help > user-guide/gen-help.txt
	if [ "${DEPLOY_URL}" != "" ]; then \
		sed -i 's#^url:.*#url: "${DEPLOY_URL}"#' _config.yaml; \
		sed -i 's#^baseurl:.*#baseurl: ""#' _config.yaml; \
	fi
	rm -rf $(dest_dir)
	bundle exec jekyll build --trace --source . --destination $(dest_dir)
	touch $(dest_dir)/.nojekyll

.PHONY: check
check: build
	bundle exec htmlproofer --disable-external --empty-alt-ignore --only-4xx --log-level :debug $(build_dir)

.PHONY: serve
serve:
	rm -rf $(dest_dir)
	bundle exec jekyll serve --source . --destination $(dest_dir) --livereload --trace
