
let tabs = {};

class Tab{
  constructor(tabId, url, time, count){
    this.tabId = tabId;
    this.url = url || null;
    this.host = url ? getHost(url) : null;
    this.startTime = time || null;
    this.count = count || 0;
    this.endTime = null;
  }
}

// 添加事件 以响应当浏览器tab更改时的状态
// 以 tab 为基础

/**
 * tab创建时记录一个空的tab元素
 */
chrome.tabs.onCreated.addListener((tab)=>{
  let _tab = new Tab(tab.id)
  tabs[tab.id] = _tab;
  console.log('tab created', tabs);
})

/**
 * 当页面更改时 响应complete事件， 根据tab来判断
 * 原本的页面(host) 增加次数
 * 新页面 更改了host  旧的页面进入保留状态，并且重置tab元素
 */
chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab)=>{
  // console.log('update触发', changeInfo);
  if(changeInfo.status === "complete"){
    let currentTab = tabs[tabId];
    if(!currentTab){
      console.log('没有找到当前tab， 可能是程序启动前创建的页面');
      return;
    }
    if(!currentTab.url){
      // tab创建后 其实是没有url的， tab的第一个页面会触发这个请求。
      currentTab.url = tab.url;
      currentTab.host = getHost(tab.url);
      currentTab.startTime = Date.now();
      currentTab.count = 1;
      console.log('tab页面创建成功', tabs, )
      return;
    }

    if(getHost(tab.url) === currentTab.host){
      // 页面更新 但是url没有改变
      ++currentTab.count;
      console.log('页面host没有变化, count++', tabs);
    } else {
      // 同一个tab下url更新了
      currentTab.endTime = Date.now();
      saveClosePage(currentTab);
      tabs[tabId] = new Tab(tab.id, tab.url, Date.now(), 1);
      console.log('tab page change', tabs)
    }
  }
})

/**
 * 当一个tab被关闭时
 * 将tab的页面内容存入数据库
 */
chrome.tabs.onRemoved.addListener((tabId,removeInfo)=>{
  // console.log('tab close');
  // console.log(tabId); 
  // console.log(removeInfo);
  console.log(removeInfo);
  let currentTab = tabs[tabId];  
  if(!currentTab){
    console.log('没有找到当前tab， 可能是程序启动前创建的页面');
    return;
  }
  // 在什么情况下等于0  
  // 只有当已经保存过了的情况下， 所以在保存无意义
  if(currentTab.count !== 0){
    currentTab.endTime = Date.now();
    saveClosePage(currentTab);  
  }
  delete tabs[tabId];
  console.log('tabs remove', tabs)
})

/**
 * 为啥要设置一个定时器
 * 因为如果是 command + q 直接关闭的 chrome整个进程 那么 onRemoved 事件是不会被触发的 所以我们找一个折中一点的方法， 定时器去轮询当前打开的页面， 
 * 有过有页面被打开了 那就保存页面打开状态以及时长信息
 */
setInterval(autoUpdateTabStatus, 30 * 1000);


/**
 * 保存关闭的页面
 */
function saveClosePage(page){
  let host = page.host;
  chrome.storage.local.get(null, function(data) {
    let items = data.items;
    if(items){
      if(hostItem = items[host]){
        console.log('找到host', host);
        items[host] =  {
          host: host,
          time: page.endTime - page.startTime + hostItem.time,
          count: hostItem.count + page.count
        }
      } else {
        items[host] =  {
          host: host,
          time: page.endTime - page.startTime,
          count: page.count
        }
        console.log('找到items， 没找到', host);
      }
    } else {
      // 需要初始化items对象
      items = {};
      items[host] =  {
        host: host,
        time: page.endTime - page.startTime,
        count: page.count
      }
      console.log('未找到items', host);
    }
    chrome.storage.local.set({items: items}, function(){
      console.log('设置成功');
    })
  });
}

/**
 * 根据url利用正则获取简单的host
 */
function getHost(url) {
  if(!url){ return null; }
  var host = url;
  var regex = /.*\:\/\/([^\/]*).*/;
  var match = url.match(regex);
  if(typeof match != "undefined" && null != match) {
    host = match[1];
  }
  return host;
}


/**
 * 自动更新当前所有tabs的状态 有需要的话保存
 */
function autoUpdateTabStatus(){
  chrome.storage.local.get(null, function(data) {
    let items = data.items;
    if(!items){
      // 从未保存过记录 抛弃这种情况
      console.log('自动保存失败, 原因:从未有过保存记录')
      return;
    }
    let flag = false;
    for(let tabId in tabs){
      let page = tabs[tabId];
      // 如果在执行自动保存过程中， tab被关闭了， 那么在这里无需保存这个tab状态
      // 因为会响应tab关闭事件的回调函数
      if(page){
        flag = true;
        let date = Date.now();
        let host = page.host;
        let item = items[host];
        if(item){
          items[host] = {
            host: host,
            time: (page.endTime || date) - page.startTime + item.time,
            count: item.count + page.count
          }
        } else {
          items[host] =  {
            host: host,
            time: (page.endTime || date)  - page.startTime,
            count: page.count
          }
        }
        page.startTime = date;
        page.count = 0;
      } 
    }
    // 如果有数据再保存 没数据就不执行
    flag && chrome.storage.local.set({items: items}, function(){
      console.log('自动更新设置成功');
    })
  })
}