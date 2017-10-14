document.addEventListener('DOMContentLoaded', () => {
  chrome.tabs.create({
    url: 'index.html'
  }, ()=>{
    console.log('xxxx');
  })
});