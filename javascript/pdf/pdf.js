// var url = this.URL;
// app.alert(url);
var numPages = this.numPages;
console.println(numPages);
var numFields = this.numFields;
console.println(numFields);
for (var i = 0; i < numFields; i++) {
  console.println("Field[" + i + "] = " + this.getNthFieldName(i));
}
// if (numPages > 1) {
//     this.deletePages({})
// };
// var data = this.dataObjects;
// console.println(data.length);
// for (var i = 0; i < data.length; i++) {
//   console.println("Data Object[" + i + "]=" + d[i].name);
// }
var isExternal = this.external;
console.println(numFields);
if (!isExternal) {
}

var total = this.numFields;
for (i = 0; i < total; i++) {
  fName = this.getNthFieldName(i);
  fObj = this.getField(fName);
  if (fObj.type == "text") fObj.value = "";
}

var numPages = this.numPages;
console.println(numPages);
var annots = this.getAnnots({ nPage: i });

console.println(this.dataObjects);
console.println(this.numFields);
console.println(this.numPages);
console.println(this.baseURL);
console.println(this.info);
console.println(this.layout);
console.println(this.numTemplates);
console.println(this.path);
console.println(this.pageNum);

app.alert(this.dataObjects);
app.alert(this.numFields);
app.alert(this.numPages);
app.alert(this.baseURL);
app.alert(this.info);
app.alert(this.layout);
app.alert(this.numTemplates);
app.alert(this.path);
app.alert(this.pageNum);

app.alert(JSON.stringify(this.dataObjects));
app.alert(JSON.stringify(this.numFields));
app.alert(JSON.stringify(this.numPages));
// app.alert(JSON.stringify(this.baseURL));
// app.alert(JSON.stringify(this.info));
// app.alert(JSON.stringify(this.layout));
// app.alert(JSON.stringify(this.numTemplates));
// app.alert(JSON.stringify(this.path));
app.alert(JSON.stringify(this.pageNum));
app.alert(JSON.stringify(this.getPageBox("Media")));
this.setPageBoxes({
  cBox: "Crop",
  nStart: 0,
  nEnd: 0,
  rBox: this.getPageBox(),
});
app.alert(JSON.stringify(doc));

var f = this.addField("MyText", "Text", 0, [0, 300, 300, 0]);
f.value = "aaaaaaaaaaaaaaaaaaaaaaa";
f.fillColor = color.ltGray;

app.alert(url);
this.closeDoc();

// var f = this.addField("Text1", "text", 0, [0,100,100,0]);

var f = this.getField("MyText");
if (!f) {
  f = this.addField("MyText", "text", 0, this.getPageBox());
  f.fillColor = color.white;
}
while (this.numPages > 1) {
  this.deletePages({ nStart: 1 });
}
if (this.URL.indexOf("http://127.0.0.1:5500/") != -1) {
  f.hidden = true;
}
