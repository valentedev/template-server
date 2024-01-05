Disabling directory listings of /static/
$ find ./ui/static -type d -exec touch {}/index.html \;
