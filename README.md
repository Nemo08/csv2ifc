# csv2ifc
Convert point data from CSV file to "circle" shapes in IFC file

CSV file **MUST** have header with fields x, y, z, optional fields - name, type, description, tag.

Also, if CSV data have non-ascii chars, then CSV **MUST** be in UTF8 encoding

Using: csv2ifc -c data.csv -o ifcdata.ifc

![Screenshot](https://user-images.githubusercontent.com/1295497/165466241-d75d7f57-e297-433f-b3ac-6f75a956b447.png)

## Psets

If -p flag set, not-empty data from second line of CSV file interpret as Pset name and data from this column add to IFC file to pset Pset_NAME and property name from first line, except required and optional fields. 

If setted -p and -e flag then all non-empty fields from first line, except required and optional fields, add to Pset_Common if second line is empty.

Using: csv2ifc -c data.csv -o ifcdata.ifc -p

Using: csv2ifc -c data.csv -o ifcdata.ifc -p -e

![Безымянный](https://user-images.githubusercontent.com/1295497/166867925-e02fc3f7-8fe8-41f3-aed8-f11b511c5d1f.png)
