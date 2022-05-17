package internal

import (
	"fmt"
	"strings"
)

var (
	IfcHeader = `
ISO-10303-21;
HEADER;
FILE_DESCRIPTION(('ViewDefinition[DesignTransferView]'), '2;1');
FILE_NAME('I:/pt_optimized.ifc','2022-04-26T13:10:55',(),(), '3.6..25428', 'Windows', 'Nobody');
FILE_SCHEMA(('IFC2X3'));
ENDSEC;
DATA;
#1= IFCAPPLICATION(#7,'0.0.220130','CSV2IFC','CSV2IFC');
#2= IFCPERSONANDORGANIZATION(#33,#32,$);
#3= IFCDIRECTION((.25,0.,0.));
#4= IFCDIRECTION((0.,0.,.25));
#5= IFCLOCALPLACEMENT(#18,#6);
#6= IFCAXIS2PLACEMENT3D(#13,#4,#3);
#7= IFCORGANIZATION('','','',(#20),(#37));
#8= IFCGEOMETRICREPRESENTATIONCONTEXT($,'Model',3,0.,#6,$);
#9= IFCGEOMETRICREPRESENTATIONSUBCONTEXT('Body','Model',*,*,*,*,#8,$,.MODEL_VIEW.,$);
#13= IFCCARTESIANPOINT((0.,0.,0.));
#14= IFCGEOMETRICREPRESENTATIONCONTEXT($,'Plan',2,0.,#38,$);
#15= IFCLOCALPLACEMENT($,#6);
#17= IFCBUILDINGSTOREY('3bWtFhoODBNOUpXCJSRsni',$,'Storey',$,$,#5,$,$,.ELEMENT.,$);
#18= IFCLOCALPLACEMENT(#15,#6);
#19= IFCSITE('2hTdvF9fv7DfUEg9JbGs8D',$,'Site',$,$,#15,$,$,.ELEMENT.,$,$,$,$,$);
#20= IFCACTORROLE(.USERDEFINED.,'',$);
#32= IFCORGANIZATION('','',$,$,$);
#33= IFCPERSON('','','',$,$,$,$,$);
#34= IFCSIUNIT(*,.LENGTHUNIT.,$,.METRE.);
#35= IFCSIUNIT(*,.AREAUNIT.,$,.SQUARE_METRE.);
#36= IFCSIUNIT(*,.VOLUMEUNIT.,$,.CUBIC_METRE.);
#37= IFCTELECOMADDRESS(.USERDEFINED.,$,'WEBPAGE',$,$,$,$,'');
#38= IFCAXIS2PLACEMENT2D(#13,#3);
#60= IFCUNITASSIGNMENT((#35,#34,#36));
#61= IFCPROJECT('2QFg6v_AvB7vKXH7g$r00z',$,'Project',$,$,$,$,(#14,#8),#60);
#63= IFCGEOMETRICREPRESENTATIONSUBCONTEXT('Box','Model',*,*,*,*,#8,$,.MODEL_VIEW.,$);
#64= IFCGEOMETRICREPRESENTATIONSUBCONTEXT('Annotation','Plan',*,*,*,*,#14,$,.PLAN_VIEW.,$);
#65= IFCRELAGGREGATES('0RFXIZw91E2ORklQ$NX2rS',$,$,$,#61,(#19));
#66= IFCRELAGGREGATES('3N$3zqYg19hOcTkmCzPGKo',$,$,$,#19,(#17));
#80= IFCFACETEDBREP(#81);
#81= IFCCLOSEDSHELL((#82));
#82= IFCFACE((#83));
#83= IFCFACEOUTERBOUND(#84,.T.);
#84= IFCPOLYLOOP((#86,#87,#88,#89,#90,#91,#92,#85));
#85= IFCCARTESIANPOINT((0.,.25,0.));
#86= IFCCARTESIANPOINT((-0.177,0.177,0.));
#87= IFCCARTESIANPOINT((-.25,0.,0.));
#88= IFCCARTESIANPOINT((-0.177,-0.177,0.));
#89= IFCCARTESIANPOINT((0.,-.25,0.));
#90= IFCCARTESIANPOINT((0.177,-0.177,0.));
#91= IFCCARTESIANPOINT((.25,0.,0.));
#92= IFCCARTESIANPOINT((0.177,0.177,0.));

`
	IfcBottom = `
ENDSEC;
END-ISO-10303-21;
`
)

func OneRecord(counter int32, x, y, z, name, itype, descr, tag string) ([]byte, int32) {
	guid, _ := NewIFCGUID()
	name, _ = Encode2HexString(name)
	descr, _ = Encode2HexString(descr)
	itype, _ = Encode2HexString(itype)
	tag, _ = Encode2HexString(tag)

	b :=
		`#` + fmt.Sprint(counter+1) + `= IFCCARTESIANPOINT((` + string(x) + `,` + string(y) + `,` + string(z) + `));
#` + fmt.Sprint(counter+2) + `= IFCPRODUCTDEFINITIONSHAPE($,$,(#` + fmt.Sprint(counter+4) + `));
#` + fmt.Sprint(counter+3) + `= IFCLOCALPLACEMENT(#5,#` + fmt.Sprint(counter+5) + `);
#` + fmt.Sprint(counter+4) + `= IFCSHAPEREPRESENTATION(#9,'Body','Brep',(#80));
#` + fmt.Sprint(counter+5) + `= IFCAXIS2PLACEMENT3D(#` + fmt.Sprint(counter+1) + `,#4,#3);
#` + fmt.Sprint(counter+6) + `= IFCBUILDINGELEMENTPROXY('` + string(guid) + `',$,'` + name + `','` + descr + `','` + itype + `',#` + fmt.Sprint(counter+3) + `,#` + fmt.Sprint(counter+2) + `,'` + tag + `',$);

`
	return []byte(b), counter + 7
}

func OnePset(counter int32, linkTo int32, name string, props map[string]string) ([]byte, int32) {
	var str string
	var nums []string

	xName, _ := Encode2HexString(name)
	nums = make([]string, 0)
	counter++

	for k, v := range props {
		xKey, _ := Encode2HexString(k)
		xValue, _ := Encode2HexString(v)
		str = str + "#" + fmt.Sprint(counter) + "= IFCPROPERTYSINGLEVALUE('" + xKey + "',$,IFCLABEL('" + xValue + "'),$);\n"
		nums = append(nums, "#"+fmt.Sprint(counter))
		counter++
	}

	guid, _ := NewIFCGUID()
	guid2, _ := NewIFCGUID()
	str = str + "#" + fmt.Sprint(counter+1) + "= IFCRELDEFINESBYPROPERTIES('" + string(guid) + "',$,$,$,(#" + fmt.Sprint(linkTo) + "),#" + fmt.Sprint(counter+2) + ");\n"
	str = str + "#" + fmt.Sprint(counter+2) + "= IFCPROPERTYSET('" + string(guid2) + "',$,'" + xName + "',$,(" + strings.Join(nums, ",") + "));\n"
	counter += 2
	str = str + "\n"

	counter++
	return []byte(str), counter
}
