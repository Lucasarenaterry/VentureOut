console.log('service worker getting out of bed and having a coffee');


self.addEventListener("install", function(event) {
   event.waitUntil(
      caches.open(CacheName)
      .then(cache => {
         return cache.addAll(CacheUrls);
      })
   );
});

self.addEventListener('activate', function(event) {
    console.log('Service worker has been activated');
    // Deletes old caches
    event.waitUntil(
        caches.keys().then( function(keys) {
            return Promise.all(keys
                .filter(key => key !== CacheName)
                .map(key => caches.delete(key))
            )
        })
    )
});


self.addEventListener('fetch', function(event) {
    // console.log('service worker fetch event', event);
    event.respondWith(
        caches.match(event.request).then( function(cacheResponse) {
            if (cacheResponse) {
                return cacheResponse;
            }
            return fetch(event.request);
        })
    );
});

const CacheName = 'VentureOut-cache-v2';
const CacheUrls = [
    './',
    '/static/css/main.css', 
    '/static/css/mdb.dark.min.css',
    '/static/css/mdb.dark.rtl.min.css',
    '/static/css/mdb.min.css',
    '/static/css/mdb.rtl.min.css',
    '/static/img/favicon.ico',
    '/static/img/logo.png',
    '/static/js/html5-qrcode.min.js',
    '/static/js/mdb.min.js',
    'manifest.webmanifest',
    'service-worker.js',
    '/static/templates/addevent.html',
    '/static/templates/ar.html',
    '/static/templates/footer.html',
    '/static/templates/header.html',
    '/static/templates/index.html',
    '/static/templates/calender.html',
    '/static/templates/infomodal.html',
    '/static/templates/map.html',
    '/static/templates/navbar.html',
    '/static/templates/navbarmobile.html',
    '/static/templates/scan.html',
    '/static/templates/scripts.html',
    '/static/templates/settings.html'
];