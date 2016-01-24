if (Meteor.isClient) {
  // counter starts at 0
  Session.setDefault('counter', 0);
  Session.setDefault('apiReturn', 'not pressed');
  Session.setDefault('freeApiReturn', 'not pressed');
  // TemplateVar.set('message', 'Not called yet');


  Template.hello.helpers({
    counter: function () {
      return Session.get('counter');
    }, 
    apiReturn: function () {
      return Session.get('apiReturn');
    }, 
    freeApiReturn: function () {
      return Session.get('freeApiReturn');
    }
  });

  Template.hello.events({
    'click button#paidButton': function () {
      console.log("requesting...");
      Session.set('apiReturn', 'Requesting it..')

      var request = new XMLHttpRequest();
      request.onreadystatechange = function() {
        if (request.readyState == XMLHttpRequest.DONE ) {
            console.log(request, request.responseText.substring(0,3));
            if (request.responseText.substring(0,3)=="404") {
              Session.set('apiReturn', 'Not found')
            }
            var response = JSON.parse(request.response)

            console.log(response)
            if (response.error) {
              Session.set('apiReturn', response.error);
            } else if (response.city) {
              Session.set('apiReturn', 'That\'s in ' + response.city)
            } else {
              Session.set('apiReturn', 'That\'s somewhere in ' + response.country_name)
            }
        }
      }

      var host = document.getElementById('host');
      request.open("GET", 'http://demo.etherapis.io:8000?ip='+host.value, true);
      request.setRequestHeader("etherapi-authorization", '{"consumer": "0x01", "provider": "0x02", "signature": "aabbff", "amount": 1000}');

      request.send();

      },

      'click button#freeButton': function () {
      console.log("requesting...");
      Session.set('freeApiReturn', 'Requesting it..')

      var request = new XMLHttpRequest();
      request.onreadystatechange = function() {
        if (request.readyState == XMLHttpRequest.DONE ) {
            console.log(request, request.responseText.substring(0,3));
            if (request.responseText.substring(0,3)=="404") {
              Session.set('freeApiReturn', 'Not found')
            }
            var response = JSON.parse(request.response)

            console.log(response)
            if (response.error) {
              Session.set('freeApiReturn', response.error);
            } else if (response.city) {
              Session.set('freeApiReturn', 'That\'s in ' + response.city)
            } else {
              Session.set('freeApiReturn', 'That\'s somewhere in ' + response.country_name)
            }
        }
      }

      var host = document.getElementById('host');
      request.open("GET", 'http://demo.etherapis.io:8000?ip='+host.value, true);
      // request.setRequestHeader("etherapi-authorization", '{"consumer": "0", "provider": "0", "signature": "0", "amount": 0}');
      // request.setRequestHeader("etherapi-authorization", '{"consumer": "0x01", "provider": "0x01", "signature": "aabbff", "amount": 1000}');

      request.send();

      }
  });
}

if (Meteor.isServer) {
  Meteor.startup(function () {
    // code to run on server at startup
  });
}

/*
$ curl -H 'etherapi-authorization: {"consumer": "0x01", "provider": "0x02", "signature": "aabbff", "amount": 1000}' http://demo.etherapis.io:8000?ip=gophergala.com

// authorization is the data content of a client-to-server payment authorization
// header, based on which the server may check for fund availability
type authorization struct {
    Consumer  string `json:"consumer"`  // API consumer authorizing the payment
    Provider  string `json:"provider"`  // API provider to which to authorize the payment to
    Amount    uint64 `json:"amount"`    // Amount of calls/data to authorize
    Signature string `json:"signature"` // Secp256k1 elliptic curve signature
}


$ curl -H 'etherapi-authorization: {"consumer": "0x01", "provider": "0x02", "signature": "aabbff", "amount": 1}' localhost:1080
{"ip":"188.24.117.6","country_code":"RO","country_name":"Romania","region_code":"CJ","region_name":"Judetul Cluj","city":"Cluj-Napoca","zip_code":"","time_zone":"Europe/Bucharest","latitude":46.7667,"longitude":23.6,"metro_code":0}

$ repeat same command
{"authorized":1,"proof":"aabbff","need":2,"error":"Not enough funds authorized"}

$ curl -H 'etherapi-authorization: {"consumer": "0x01", "provider": "0x02", "signature": "aabbff", "amount": 2}' localhost:1080
{"ip":"188.24.117.6","country_code":"RO","country_name":"Romania","region_code":"CJ","region_name":"Judetul Cluj","city":"Cluj-Napoca","zip_code":"","time_zone":"Europe/Bucharest","latitude":46.7667,"longitude":23.6,"metro_code":0}

$ repeat same command
{"authorized":2,"proof":"aabbff","need":3,"error":"Not enough funds authorized"}


*/



