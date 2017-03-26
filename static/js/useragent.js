/**
 * chrome
 * windows  Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.110 Safari/537.36
 * mac      Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.1 Safari/537.36
 * iphone   Mozilla/5.0 (iPhone; CPU iPhone OS 8_1 like Mac OS X)
 * android  Mozilla/5.0 (Linux; Android 5.1; XT1040 Build/LPB23.13-35) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.84 Mobile Safari/537.36
 *
 * edge
 * win10    Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15061
 *
 * ie11
 * win10    Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko
 */

const UserAgent = {
    os: /^Mozilla\/5\.0\s\((.+?)\)/,//->Windows NT 10.0 提取括号里第一个分号
    unknownOs: /(.+?);/,
    windowsName: {
        '5.1': ["Windows XP"],
        '5.2': ["Windows Server 2003", "Windows Server 2003 R2"],
        '6.0': ["Windows Vista", "Windows Server 2008"],
        '6.1': ["Windows 7", "Windows Server 2008 R2"],
        '6.2': ["Windows 8.0", "Windows Server 2012"],
        '6.3': ["Windows 8.1", "Windows Server 2012 R2"],
        '10.0': ["Windows 10"],
    },

    //获得os type和os版本
    getOsType: function (s) {
        let version = 'unknown';
        let ostype = 'unknown';
        let os = s.match(this.os);//匹配第一个括号里面的Windows NT 10.0; Win64; x64
        if (os) {
            console.log(os);
            let t = os[1];
            let v;
            if (t.startsWith("Windows")) {//windows 设备
                if (v = t.match(/Windows\sNT\s(\d+\.\d+)/)) {
                    version = v[1];
                }
                let n = this.windowsName[version];
                ostype = n ? n[0] : "Windows";
            } else if (t.match(/OS (.*) like Mac OS X/)) {//ios设备
                ostype = t.match(/(.+?);/)[1];//iPhone iPad iPod
                if (v = t.match(/CPU\siPhone\sOS\s([\d_]+)\s/)) {
                    version = v[1].replace('_', '.');
                }
            } else if (t.startsWith("Linux")) {//linux 设备
                if (t.includes('Android')) {
                    ostype = "Android";
                    if (v = t.match(/Android\s([\d.]+);/)) {
                        version = v[1];
                    }
                } else if (t.includes('Ubuntu')) {
                    ostype = 'Ubuntu';
                    if (v = /Ubuntu\/([0-9.]*)/.exec(t)) {
                        version = v[1];
                    }
                } else if (t.includes('Debian')) {
                    ostype = 'Debian';
                } else if (t.includes('CentOS')) {
                    ostype = 'CentOS';
                    if (v = /CentOS\/[0-9\.\-]+el([0-9_]+)/.exec(t)) {
                        version = v[1].replace(/_/g, '.');
                    }
                } else if (t.includes('SUSE')) {
                    ostype = 'SUSE';
                } else if (t.includes('Fedora')) {
                    ostype = 'Fedora';
                    if (v = /Fedora\/[0-9\.\-]+fc([0-9]+)/.exec(t)) {
                        version = v[1];
                    }
                } else if (t.includes('Gentoo')) {
                    ostype = 'Gentoo';
                } else if (t.includes('Kubuntu')) {
                    ostype = 'Kubuntu';
                } else if (t.includes('Slackware')) {
                    ostype = 'Slackware';
                } else if (t.includes('Red Hat')) {
                    ostype = 'Red Hat';
                    if (v = /Red Hat[^\/]*\/[0-9\.\-]+el([0-9_]+)/.exec(t)) {
                        version = v[1].replace(/_/g, '.')
                    }
                } else if (t.includes('Mageia')) {
                    ostype = 'Mageia';
                    if (v = /Mageia\/[0-9\.\-]+mga([0-9]+)/.exec(t)) {
                        version = v[1];
                    }
                } else if (t.includes('Mandriva Linux')) {
                    ostype = 'Mandriva';
                    if (v = /Mandriva Linux\/[0-9\.\-]+mdv([0-9]+)/.exec(t)) {
                        version = v[1];
                    }
                }
            } else if (t.startsWith("Macintosh")) {
                ostype = "Mac";
                if (v = t.match(/Mac\sOS\sX\s([\d_]+)/)) {
                    version = v[1].replace(/_/g, '.')
                }
            } else if (t.startsWith('Unix')) {
                ostype = "Unix";
            } else if (t.startsWith('FreeBSD')) {
                ostype = "FreeBSD";
            } else if (t.startsWith('OpenBSD')) {
                ostype = "OpenBSD";
            } else if (t.startsWith('NetBSD')) {
                ostype = "NetBSD";
            } else if (t.startsWith('Solaris')) {
                ostype = "Solaris";
            } else if (t.includes('Tizen')) {
                ostype = 'Tizen';
                if (v = /Tizen[\/ ]([0-9.]*)/.exec(t)) {
                    version = v[1];
                }
            } else {
                let t2 = t.match(this.unknownOs);
                ostype = t2 ? t2[1] : "unknown os";
                version = "unknown";
            }
        }

        return {
            os: ostype,
            version: version
        }
    },


};