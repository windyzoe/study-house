# \*\*-地图找房字段分析

- border:经纬度坐标数据,画出地图上一个块
- longitude:经度
- latitude:纬度
- bubble:一个地图块,可能是小区\板块\区

# 表结构设计

## 房子 house (因为 spider,设计了一些冗余字段)

| 字段       | 描述     |
| ---------- | -------- |
| Name       | 名称     |
| Building   | 小区     |
| Region     | 板块     |
| District   | 市区     |
| TotalPrice | 总价     |
| UnitPrice  | 单价     |
| Time       | 修建时间 |
| Area       | 面积     |
| FloorCount | 楼层     |

## 小区 building

| 字段          | 描述                      |
| ------------- | ------------------------- |
| Name          | 名称                      |
| Region        | 板块                      |
| District      | 市区                      |
| Time          | 修建时间                  |
| Decription    | 备注                      |
| BuildingCount | 楼数                      |
| HouseCount    | 房屋总数                  |
| Alias         | 别名,小区可能有不同的街号 |

## 学校 school

| 字段     | 描述     |
| -------- | -------- |
| Name     | 名称     |
| Region   | 板块     |
| District | 市区     |
| Size     | 小初中   |
| Awesome  | 牛逼程度 |

## 学校-小区关联表 rel_school_building

| 字段       | 描述         |
| ---------- | ------------ |
| Id         | id           |
| Building   | 小区         |
| School     | 学校         |
| Year       | 几几年的规定 |
| Decription | 备注         |

# 数据清洗

## 小区-学校映射

- 小区结尾 "弄"
- "弄"前的数字为阿拉伯文,其余数字为中文
- $为占位,填充集中类型如下
  > [1,2,3] : 枚举
  > {1,100}: 范围
  > (单)(双): 单双号
  > :>100: 范围
  > :<100: 范围
  > &&: 且关系
  > e1,2,3e: 除了
  > 或关系单独编辑行
