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
    "/", 
    "app.js", 
    "styles.css", 
    "logo.svg"
];