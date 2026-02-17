const { fas } = require('@fortawesome/free-solid-svg-icons');

// 检查特定图标的实际名称
const testKeys = ['faDiceD6', 'faDiceD20', 'faArrowUp19', 'faArrowUp91', 'faList12'];
console.log('Actual icon names:');
testKeys.forEach(k => {
  const icon = fas[k];
  if (icon) {
    console.log('  ' + k + ' -> iconName: ' + icon.iconName);
  }
});
