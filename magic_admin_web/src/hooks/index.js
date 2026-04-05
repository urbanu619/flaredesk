// 递归菜单结构
export const transformMenu = (menuItems, transformFn) => {
  return menuItems.map(item => {
    const transformedItem = transformFn(item);
    if (item.children && item.children.length) {
      transformedItem.children = transformMenu(item.children, transformFn);
    }
    return transformedItem;
  });
}


//一维 数组按属性拆分为二维数组
export const arrToArr = (arr, prop) => {
  const result = [];
  const rArr = []

  for (const item of arr) {
    const key = item[prop];
    if (!result[key]) {
      result[key] = [];
    }
    result[key].push(item);
  }

  Object.keys(result).forEach((item, index) => {
    console.log(result[item])
    let item_0 = result[item].filter(a => {
      return a.lang === 'zh'
    })[0];
    console.log(item_0)

    let obj = {
      index: index + 1,
      id: item_0?.id || "",
      title: item_0?.answer || "",
      questionNo: result[item][0].questionNo,
      question: item_0?.question,
      children: result[item],
    }
    rArr.push(obj);
  })
  return rArr;
}

import dayjs from "dayjs";

export const dataTo_YMD = (data) => {
  if (data) {
    return dayjs(data).format("YYYY-MM-DD")
  } else {
    return ""
  }


}
