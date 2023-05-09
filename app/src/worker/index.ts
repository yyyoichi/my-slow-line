self.addEventListener('install', () => {
  console.log('install sw');
});
self.addEventListener('push', (event) => {
  console.log(event);
});
self.addEventListener('load', (event) => {
  console.log('load', event);
});
