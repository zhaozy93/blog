chrome.storage.local.get(null, (data)=>{
  drawCount(filterHost(data.items, 'count'));
  drawTime(filterHost(data.items, 'time'));
  console.log('您的所有浏览记录');
  console.log('但是我猜你看不到这里');
  console.log('*****************************');
  console.log(data.items);
  console.log('*****************************');
})

function filterHost(items, property, number = 50){
  let result = [];
  let min = 0, max = 0;

  for(let host in items){
    let page = items[host];
    let value = page[property];
    if(page.count >0){
        if(value > max){
            result.push(page);
            max = value;
            if(result.length > number){
              result.shift();
              min = result[0][property];
            }
            if(result.length === 1){
              min = value;
            }
          } else if(value < min){
            if(result.length < number){
              result.unshift(page);
              min = value;
            }
          } else {
            for(let i = 1; i < result.length; i++){
              if(value >= result[i - 1][property] && value <= result[i][property]){
                result.splice(i, 0 , page);
                if(result.length > number){
                  result.shift();
                  min = result[0][property];
                }
                break;
              }
            }
          }
    }
    
  }
  return result;
}

function drawCount(data){
  if(!data || !data.length){
    document.getElementById('counts').innerHTML = '<h1>暂无数据</h1>'
    return;
  }
  let options = {
    chart: {
        type: 'bar'
    },
    title: {
        text: '访问量前50名'
    },
    xAxis: {
        // categories: ['Africa', 'America', 'Asia', 'Europe', 'Oceania'],
        title: {
            text: null
        }
    },
    yAxis: {
        min: 0,
        title: {
            text: '次数',
            align: 'high'
        },
        labels: {
            overflow: 'justify'
        }
    },
    tooltip: {
        valueSuffix: '次'
    },
    plotOptions: {
        bar: {
            dataLabels: {
                enabled: true
            }
        }
    },
    // legend: {
    //     layout: 'vertical',
    //     align: 'right',
    //     verticalAlign: 'top',
    //     x: -40,
    //     y: 80,
    //     floating: true,
    //     borderWidth: 1,
    //     backgroundColor: ((Highcharts.theme && Highcharts.theme.legendBackgroundColor) || '#FFFFFF'),
    //     shadow: true
    // },
    credits: {
        enabled: false
    },
    series: [{
        name: '访问次数',
        // data: [107, 31, 635, 203, 2]
    }]
  }
  options.xAxis.categories = data.map(item=>{
    return item.host;
  })
  options.series[0].data = data.map(item=>{
    return item.count;
  })
  Highcharts.chart('counts', options);
}

function drawTime(data){
  if(!data || !data.length){
    document.getElementById('times').innerHTML = '<h1>暂无数据</h1>'
    return;
  }
  let options = {
    chart: {
        type: 'bar'
    },
    title: {
        text: '访问总时长前50名'
    },
    xAxis: {
        title: {
            text: null
        }
    },
    yAxis: {
        min: 0,
        title: {
            text: '时长/分',
            align: 'high'
        },
        labels: {
            overflow: 'justify'
        }
    },
    tooltip: {
        valueSuffix: '分'
    },
    plotOptions: {
        bar: {
            dataLabels: {
                enabled: true
            }
        }
    },
    legend: {
        layout: 'vertical',
        align: 'right',
        verticalAlign: 'top',
        x: -40,
        y: 80,
        floating: true,
        borderWidth: 1,
        backgroundColor: ((Highcharts.theme && Highcharts.theme.legendBackgroundColor) || '#FFFFFF'),
        shadow: true
    },
    credits: {
        enabled: false
    },
    series: [{
        name: '访问总时长'
    },{
      name: '每页面平均时长'
    }]
  }
  options.xAxis.categories = data.map(item=>{
    return item.host;
  })
  options.series[0].data = data.map(item=>{
    return +(item.time / 1000 / 60).toFixed(2);
  });
  options.series[1].data = data.map(item=>{
    return +(item.time / item.count / 1000 / 60).toFixed(2);
  })
  Highcharts.chart('times', options);
}


document.getElementById('clear').addEventListener('click', ()=>{
    chrome.storage.local.clear(()=>{
        alert('清除成功');
    })
}, false)

