console.log('service worker getting out of bed and having a coffee');


self.addEventListener("install", function(event) {
   event.waitUntil(
      caches.open("VentureOut-cache")
      .then(cache => {
         return cache.addAll(CacheUrls);
      })
   );
});

self.addEventListener('activate', function(event) {
    console.log('Service worker has been activated')
});


const CacheUrls = [
    './',
    '/static/css/main.css', 
    '/static/css/mdb.dark.min.css',
    '/static/css/mdb.dark.rtl.min.css',
    '/static/css/mdb.min.css',
    '/static/css/mdb.rtl.min.css',
    '/static/img/favicon.ico',
    '/static/img/logo.png',
    '/static/img/maskable_icon_x72.png',
    '/static/img/maskable_icon_x96.png',
    '/static/img/maskable_icon_x128.png',
    '/static/img/maskable_icon_x192.png',
    '/static/js/html5-qrcode.min.js',
    '/static/js/mdb.min.js',
    'manifest.webmanifest',
    'service-worker.js',
    '/static/HWUevents.geojson',
    '/static/templates/addevent.html',
    '/static/templates/footer.html',
    '/static/templates/header.html',
    '/static/templates/index.html',
    '/static/templates/infomodal.html',
    '/static/templates/map.html',
    '/static/templates/navbar.html',
    '/static/templates/navbarmobile.html',
    '/static/templates/scan.html',
    '/static/templates/scripts.html',
    '/static/templates/settings.html'
];