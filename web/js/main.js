angular.module('vpnrApp', ["ngFlash"])
.factory('Base64', function() {
    var keyStr = 'ABCDEFGHIJKLMNOP' +
        'QRSTUVWXYZabcdef' +
        'ghijklmnopqrstuv' +
        'wxyz0123456789+/' +
        '=';
    return {
        encode: function (input) {
            var output = "";
            var chr1, chr2, chr3 = "";
            var enc1, enc2, enc3, enc4 = "";
            var i = 0;

            do {
                chr1 = input.charCodeAt(i++);
                chr2 = input.charCodeAt(i++);
                chr3 = input.charCodeAt(i++);

                enc1 = chr1 >> 2;
                enc2 = ((chr1 & 3) << 4) | (chr2 >> 4);
                enc3 = ((chr2 & 15) << 2) | (chr3 >> 6);
                enc4 = chr3 & 63;

                if (isNaN(chr2)) {
                    enc3 = enc4 = 64;
                } else if (isNaN(chr3)) {
                    enc4 = 64;
                }

                output = output +
                    keyStr.charAt(enc1) +
                    keyStr.charAt(enc2) +
                    keyStr.charAt(enc3) +
                    keyStr.charAt(enc4);
                chr1 = chr2 = chr3 = "";
                enc1 = enc2 = enc3 = enc4 = "";
            } while (i < input.length);

            return output;
        },

        decode: function (input) {
            var output = "";
            var chr1, chr2, chr3 = "";
            var enc1, enc2, enc3, enc4 = "";
            var i = 0;

            // remove all characters that are not A-Z, a-z, 0-9, +, /, or =
            var base64test = /[^A-Za-z0-9\+\/\=]/g;
            if (base64test.exec(input)) {
                alert("There were invalid base64 characters in the input text.\n" +
                        "Valid base64 characters are A-Z, a-z, 0-9, '+', '/',and '='\n" +
                        "Expect errors in decoding.");
            }
            input = input.replace(/[^A-Za-z0-9\+\/\=]/g, "");

            do {
                enc1 = keyStr.indexOf(input.charAt(i++));
                enc2 = keyStr.indexOf(input.charAt(i++));
                enc3 = keyStr.indexOf(input.charAt(i++));
                enc4 = keyStr.indexOf(input.charAt(i++));

                chr1 = (enc1 << 2) | (enc2 >> 4);
                chr2 = ((enc2 & 15) << 4) | (enc3 >> 2);
                chr3 = ((enc3 & 3) << 6) | enc4;

                output = output + String.fromCharCode(chr1);

                if (enc3 != 64) {
                    output = output + String.fromCharCode(chr2);
                }
                if (enc4 != 64) {
                    output = output + String.fromCharCode(chr3);
                }

                chr1 = chr2 = chr3 = "";
                enc1 = enc2 = enc3 = enc4 = "";

            } while (i < input.length);

            return output;
        }
    };
})
.controller('RouteController', function(Base64,$location, $scope, $http, Flash) {
    var routeList = this;
    var endpoint = "/api";
    routeList.myRoute = null;
    routeList.routes = [];
    routeList.tables = [];
    var init = function() {
        var encoded = Base64.encode($location.hash());
        $http.defaults.headers.common['Authorization'] = 'Bearer ' + encoded;
    };
    init();
    var load = function() {
        $http.get(endpoint+"/tables").success(function(data){
            routeList.tables = data.data;
        });
        $http.get(endpoint+"/routes").success(function(data){
            preprocessData(data); 
        });
    };
    var preprocessData = function(data) {
        if (!data.data) {
            return
        }
        var ip = data["request-ip"];
        routeList.routes = new Array();
        for (i=0;i<data.data.length; i++) {
            if (data.data[i].ip == ip) {
                routeList.myRoute = data.data[i]; 
            }else{
                routeList.routes.push(data.data[i]);
            }
        }
    };
    load();
    routeList.setRoute = function(ip, table) {
        $http.post(endpoint+"/routes", {data:{ip: ip, table: table}}).success(function(data){
            load();

        }).error(function(data){
            Flash.create('danger', "<strong>Permission denied</strong>", 2000, {class: 'alert alert-danger navbar-alert', id:'navbar-alert'}, false); 
        });

    };
    routeList.tableClasses = [
        "btn-default",
        "btn-primary",
        "btn-success",
        "btn-info",
        "btn-warning",
        "btn-danger",
    ];

    routeList.tableByName = function(name) {
        for (i=0;i<routeList.tables.length;i++) {
            if (routeList.tables[i].name == name) {
                return routeList.tables[i];
            }
        }
        return null;
    };
    routeList.tableClass = function(name) {
        for (i=0;i<routeList.tables.length;i++) {
            if (routeList.tables[i].name == name) {
                return routeList.tableClasses[i];
            }
        }
        return routeList.tableClasses[0];
    };

    routeList.missingTables = function(name) {
        tables = new Array(); 
        for (i=0; i<routeList.tables.length;i++) {
            t=routeList.tables[i];
            if (t.name != name) {
                tables.push(t);
            }
        }
        return tables;
    };
});

