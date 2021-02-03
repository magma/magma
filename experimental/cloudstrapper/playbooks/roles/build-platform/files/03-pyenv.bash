curl https://pyenv.run | bash

echo "export PATH=$HOME/.pyenv/bin:"'$PATH' >> ~/.bash_profile
echo -e 'if command -v pyenv 1>/dev/null 2>&1; then\n  eval "$(pyenv init -)"\nfi' >> ~/.bash_profile

sudo ln -s ~/.pyenv/bin/pyenv /usr/local/bin

pyenv install 3.7.3
pyenv global 3.7.3
