/* eslint-disable */
import fetch from 'node-fetch';
import express from 'express';
import morgan from 'morgan';

import Pyroscope from 'pyroscope';

const port = process.env['PORT'] || 3000;

const region = process.env['REGION'] || 'default';

const app = express();
app.use(morgan('dev'));

app.get('/', (req, res) => {
  res.send('Available routes are: /bike, /car, /scooter');
});

function tagSpan(name) {
  return eval('obj = { ' + name + ': function (fn) { return fn } }.' + name);
}


// Object.defineProperty(adder, 'name', {
//   value: 'subber',
//   writable: false
// });
//(adder as any).name = "subber";
// (adder as any).displayName = "subber333";
const subber = tagSpan("subber4545");
console.log(subber.name);
const subberFn = subber(function (a,b ) {
  if ( b.toString()[1] === '3' ) {
    return a +  b * 2;
}
  return a + b;
});
console.log(subberFn.name);

const genericSearchHandler = (p: number) => (req: any, res: any) => {
  const time = +new Date() + p * 1000;
  let i = 0;
  while (+new Date() < time) {
    i = subberFn(i, Math.random());
  }
  res.send('Vehicle found');
};

app.get('/bike', function bikeSearchHandler(req, res) {
  return genericSearchHandler(0.2)(req, res);
});
app.get('/car', function carSearchHandler(req, res) {
  return genericSearchHandler(1)(req, res);
});
app.get('/scooter', function scooterSearchHandler(req, res) {
  return genericSearchHandler(0.5)(req, res);
});

setInterval(() => {
  fetch(`http://localhost:${port}/car`);
}, 1800);

setInterval(() => {
  fetch(`http://localhost:${port}/bike`);
}, 633);

setInterval(() => {
  fetch(`http://localhost:${port}/scooter`);
}, 1000);

Pyroscope.init();


app.listen(port, () => {
  console.log(
    `Server has started on port ${port}, use http://localhost:${port}`
  );
});
