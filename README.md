# csv2ifc
Convert point data from CSV file to "circle" shapes in IFC file

CSV file **MUST** have header with fields x, y, z, optional fields - name, type, description, tag

Also, if CSV data have non-ascii chars, then CSV **MUST** be in UTF8 encoding


Using: csv2ifc -c data.csv -o ifcdata.out

![Screenshot](https://user-images.githubusercontent.com/1295497/165466241-d75d7f57-e297-433f-b3ac-6f75a956b447.png)
