# beanstalk-migrator [![Build Status](https://travis-ci.org/lbarulski/beanstalk-migrator.svg?branch=master)](https://travis-ci.org/lbarulski/beanstalk-migrator)

## Usage
[![asciicast](https://asciinema.org/a/ybzH40Vpi90flGjLdUr7fzaiK.png)](https://asciinema.org/a/ybzH40Vpi90flGjLdUr7fzaiK)

## Caveats
All _buried_ jobs will be moved as _ready_ - In a nutshell, [beanstalk protocol](https://github.com/kr/beanstalkd/blob/master/doc/protocol.txt#L81) does not allow to burry task in a simple manner

## FAQ
### Will it work under load?
Yes, it will. **HOWEVER** it might left some jobs on source instance and migration will take **significantly** more time.
