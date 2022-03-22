import pprof from 'pprof';

const obj = {
    name: function() {
        return 0;
    }
}

pprof.time.profile({
    lineNumbers: true,
    durationMillis: 1000,
}).then(d =>  {
    if ( d.stringTable.indexOf('taggedFunction') !== -1 ) {
        console.log("YAY");
        console.log(d.stringTable);

    } else { 
        console.log("NAH");
        console.log(d.stringTable);
    }
}
)
const globalObject = {};
setTimeout(() => {
    
    const tagSpan = (name) => {
        globalObject[name] = new Function("fn", '"use strict";return function ' + name + '(fn) { return fn }');
        return globalObject[name];
    }

    function doStuff(callback) {
        return callback
    }




    
    const fn = tagSpan("taggedFunction")();
    const tfn = fn(function() {
        var q = 24;
        for( var j = 0; j < 10000; j++ ) {
            q = j + q*Math.sin(Math.random()*j)*doStuff(function doComplexMath(a) { return a*Math.random()})(10);
        }
        return q;
    });

    console.log(fn.name);
    console.log(tfn.name);
    function sampleName() {
        var q = 24;
        for( var j = 0; j < 10000; j++ ) {
            q = j + q*Math.sin(Math.random()*j);
        }
        return q;
    }
    console.log("Strarting count");
    var q = 0;
    for ( var i = 0; i < 10000; i++ ) {
        q += sampleName();
        q += q * tfn();
    
    }    
}, 0)

