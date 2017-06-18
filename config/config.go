package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	Name         string              `yaml:"name"`
	Colors       bool                `yaml:"colors"`
	Debug        bool                `yaml:"debug"`
	MaxSize      int64               `yaml:maxsize`
	EveryXstatus int                 `yaml:everyXstatus`
	GlobList     []string            `yaml:"globlist"`
	SearchList   []string            `yaml:"searchlist"`
	IgnoreList   []string            `yaml:"ignorelist"`
	BlackList    []string            `yaml:"blacklist"`
	Commands     map[string][]string `yaml:"commands"`
}

var Testing = `
name: 'Murmamau'
colors: false
debug: false
everyxstatus: 100000
maxsize: 10485760 #  1024 * 1024 * 10   => 10 MB
searchlist:
  - /etc/passwd
  - /etc/fstab
  - /etc/shadow
  - /etc/shadow~
  - /etc/shadow-
  - /etc/master.passwd
  - /.secure/etc/passwd
  - /etc/spwd.db
  - /tcb/files/auth
  - /etc/udb
  - /etc/secrets
  - /etc/d_passwd
  - /etc/secrets
  - /etc/opasswd
  - /etc/hardened-shadow
  - /root/passwords.txt
  - /etc/sudoers
  - /etc/hosts
  - /etc/grsec/pw
  - /etc/ssh/sshd_config
  - /root/anaconda-ks.cfg
  - /root/ks.cfg
  - /etc/wpa_supplicant.conf
  - /usr/sbin/john.pot
  - /etc/tripwire/site.key
  - /etc/nagios/nsca.cfg
  - /etc/nagios/send_nsca.cfg
  - /etc/aiccu.conf
  - /etc/cntlm.conf
  - /etc/conf.d/hsqldb
  - /etc/conf.d/openconnect
  - /etc/mysql/mysqlaccess.conf
  - /etc/nessus/nessusd.conf
  - /etc/nikto/nikto.conf
  - /etc/postfix/saslpass
  - /var/lib/samba/private/smbpasswd
  - /etc/screenrc
  - /etc/sysconfig/rhn/osad-auth.conf
  - /etc/sysconfig/rhn/osad.conf
  - /etc/sysconfig/rhn/rhncfg-client.conf
  - /etc/sysconfig/rhn/up2date
  - /tmp/krb5.keytab
  - /proc/config.gz
  - /var/log/auth.log
  - /var/log/secure.log
  - ~/.bash_history
  - ~/.sh_history
  - ~/.history
  - ~/.zsh_history
  - ~/.csh_history
  - ~/.tcsh_history
  - ~/.ksh_history
  - ~/.ash_history
  - ~/.php_history
  - ~/.mysql_history
  - ~/.sqlite_history
  - ~/.psql_history
  - ~/.mc/history
  - ~/.atftp_history
  - ~/.irb_history
  - ~/.scapy_history
  - ~/.sqlplus_history
  - ~/.cvspass
  - ~/.john/john.pot
  - ~/.ssh/config
  - ~/.netrc
  - ~/.rhosts
  - ~/.shosts
  - ~/.my.cnf
  - ~/.armitage.prop
  - ~/.java/.userPrefs/burp/prefs.xml
  - ~/.ZAP/config.xml
  - ~/.filezilla/filezilla.xml
  - ~/.filezilla/recentservers.xml
  - ~/.ncftp/firewall
  - ~/.gftp/gftprc
  - ~/.gftp/bookmarks
  - ~/.remmina/remmina.pref
  - ~/.subversion/config
  - ~/.subversion/servers
  - ~/.config/gmpc/profiles.cfg
  - ~/.config/mc/ini
  - ~/.config/vlc/vlcrc
  - ~/.rootdpass
  - ~/.gitconfig
  - ~/.gnupg/secring.gpg
  - ~/.nessusrc
  - ~/.smb.cnf
  - ~/.muttrc
  - ~/.msf4/config
  - ~/.msf4/logs/console.log
  - ~/.msf4/logs/framework.log
  - ~/.msf4/logs/production.log
  - ~/.msf4/history
  - ~/.msf3/config
  - ~/.msf5/config
  - ~/.ssh/authorized_keys
  - ~/.aws/credentials
  - ~/.aws/config
  - /etc/aws/config
  - /etc/aws/credentials
  - ~/.Rhistory
  - ~/.Rapp.history
  - ~/.muttrc
  - ~/.netrc
  - .bash_history
  - .sh_history
  - .zsh_history
  - .history
  - .mysql_history
  - .psql_history
  - accounts.xml
  - id_dsa
  - id_rsa
  - id_ecdsa
  - id_ed25519
  - identity
  - connect.inc
  - default.pass
  - .htaccess
  - cert8.db
  - users.xml
  - shadow
  - passwd
  - master.passwd
  - john.pot
  - database.yml
  - secrets.yml
  - secret_token.rb
  - bower.json
  - config.json
  - keys.db
  - .dbeaver-data-sources.xml
  - .s3cfg
  - .npmrc
  - .travis.yml
  - .htpasswd
  - config.gypi
  - config.php
  - database.php
  - config.inc
  - database.inc
  - parameters.yml
  - parameters.ini
  - wp-config.php
  - settings.php
  - settings.inc
  - configuration.php
  - web.xml
  - web.config
  - secret_token.rb
  - omniauth.rb
  - carrierwave.rb
  - credentials.xml
  - schema.rb
  - LocalSettings.php
  - Dockerfile
  - credentials
  - docker-compose.yml
globlist:
  - "*password*"
  - ".ssh/*"
  - "/etc/*"
  - "*.*history"
  - "*_rsa"
  - "*.dsa"
  - "*__ed25519"
  - "*_ecdsa"
  - "*.pem"
  - "*.ppk"
  - "*.pkcs12"
  - "*.p12"
  - "*.ovpn"
  - "*.kdb"
  - "*.agilekeychain"
  - "*.keychain"
  - "*.keystore"
  - "*.keyring"
  - "*.kwallet"
  - "*.tblk"
  - "*credentials*"
  - "private.*key"
ignorelist:
  - .mp4
  - .avi
  - .mp3
  - .png
  - .jpg
  - .jpeg
  - .gif
  - .mp3
  - .mpa
  - .wav
  - .wma
  - .asf
  - .mov
  - .flv
  - .asf
  - .vob
  - .bmp
  - .tiff
  - .tif
  - .js
  - .css
  - .iso
blacklist:
  - /etc/XXXXXXXXXXX
commands:
  linux:
    - uname -a
    - id
    - pwd
    - w
    - who
    - last
    - lastlog
    - ifconfig -a
    - netstat -rn
    - ps auxef
    - lsmod
    - rpm -qa
    - glsa-check -v -t all
    - mount
    - netstat -nalp
    - lsof
    - iptables -L -n
    - hostname
    - cat /etc/shadow
    - cat /etc/master.passwd
    - sudo -S -l
    - sudo -V
    - env
    - export
    - set
    - echo $PATH
    - ifconfig -a
    - arp -a
    - ip addr
    - netstat -rn
    - route
    - netstat -antp
    - netstat -anup
    - mysqladmin -uroot -proot version
    - mysqladmin -uroot version
    - psql -U postgres template0 -c 'select version()'
    - psql -U postgres template1 -c 'select version()'
    - psql -U pgsql template0 -c 'select version()'
    - psql -U pgsql template1 -c 'select version()'
    - apache2ctl -M
    - dpkg --list
    - head /var/mail/root
    - cat /proc/self/cgroup
    - ac
    - lastcomm
    - sa
    - aureport --avc
    - aureport --auth
    - aureport --comm
    - aureport --tty
    - aureport --terminal
    - aureport --executable
    - xm list
    - gradm -C
    - sestatus -v
    - apparmor_status
  freebsd:
    - uname -a
    - id
`

/*
  read a YAML config and return config.Configuration struc
*/
func ReadConfig(conf string) (Configuration, error) {
	configuration := Configuration{}
	// Read the YAML data into a struct instance called person
	err := yaml.Unmarshal([]byte(conf), &configuration)
	if err != nil {
		fmt.Println("Error in YAML")
		return Configuration{}, err
	}
	return configuration, nil
}
