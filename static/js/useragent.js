/**
 * chrome
 *Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.110 Safari/537.36
 * Windows NT
 * Macintosh
 *
 * edge
 * Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36 Edge/15.15061
 *
 * ie 11
 * Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; .NET4.0C; .NET4.0E; rv:11.0) like Gecko
 *
 * iphone Safari
 * Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/603.1.23 (KHTML, like Gecko) Version/10.0 Mobile/14E5239e Safari/602.1
 *
 * ios chrome
 * Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.75 Mobile/14E5239e Safari/602.1
 *
 * android chrome
 * Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19
 * */


const UserAgent = {
    ntVersion: /Windows\sNT\s(\d+\.\d+)/,
    iosVersion: /CPU\siPhone\sOS\s([\d_]+)\s/,
    androidVersion: /Android\s([\d.]+);/,
    os: /^Mozilla\/5\.0\s\(([0-9a-zA-Z\s.]*);/,//->Windows NT 10.0 提取括号里第一个分号
    windowsName: {
        '5.1': "Windows XP",
        '5.2': ["Windows Server 2003", "Windows Server 2003 R2"],
        '6.0': ["Windows Vista", "Windows Server 2008"],
        '6.1': ["Windows 7", "Windows Server 2008 R2"],
        '6.2': ["Windows 8.0", "Windows Server 2012"],
        '6.3': ["Windows 8.1", "Windows Server 2012 R2"],
        '10.0': "Windows 10",
    },

    getOsType: function (s) {
        let version = 'unknown';
        let ostype = 'unknown';
        let os = s.match(this.os);
        if (os) {
            if (os[1].startsWith("Windows")) {
                version = os[1].match(this.ntVersion)[1];
                ostype = this.windowsName[version] ? this.windowsName[version] : 'unknown';
            } else if (os[1].startsWith("iPhone")) {
                version = s.match(this.iosVersion)[1].replace('_', '.');
                ostype = "iPhone";
            } else if (os[1].startsWith("Linux")) {
                version = s.match(this.androidVersion)[1].replace('_', '.');
                ostype = "Android";
            }
        }

        return {
            ostype: ostype,
            version: version
        }
    },


};