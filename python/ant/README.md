# Ant

```sh
mkdir -p ~/dev/all/python/ant
code $_
virtualenv -p python3 venv
mkdir .vscode
touch .vscode/settings.json
echo '{ "python.venvPath": "./venv" }' >> .vscode/settings.json
. venv/bin/activate
```