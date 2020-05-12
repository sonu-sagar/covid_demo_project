import React from 'react';
// import ReactDOM from 'react-dom';
import FusionCharts from "fusioncharts";
import charts from "fusioncharts/fusioncharts.charts";
import ReactFusioncharts from "react-fusioncharts";

charts(FusionCharts);

class MyComponent2 extends React.Component {
  constructor(props){
    super(props);
  }

    render() {
    var dataSource;
    var x=this.props.data;
    var lab=[];
    var currcases =[];
    var currrecovered =[];
    var currdeaths =[];
    var i;
    for(i=0;i<x.length;i++){
        lab.push({label:x[i]["Stats"]["Date"]});
        currcases.push({value:x[i]["Stats"]["CurrentDayCases"]});
        currrecovered.push({value:x[i]["Stats"]["CurrentDayRecovered"]});
        currdeaths.push({value:x[i]["Stats"]["CurrentDayDeaths"]});
    }
        dataSource = {
            chart: {
                theme: "fusion",
                caption: "Timeline of covid19",
                xAxisname: "Date",
                yAxisName: "Total Cases",
                numberPrefix: "",
                plotFillAlpha: "80",
                divLineIsDashed: "1",
                divLineDashLen: "1",
                divLineGapLen: "1"
              },
              categories: [{
                category: lab
              }],
              dataset: [{
                seriesname: "Confirmed Case",
                data: currcases
              }, {
                seriesname: "Recovered Case",
                data: currrecovered
              },{
                seriesname: "Deaths",
                data: currdeaths
              }],
            }
        

    return (
        <div>
            <ReactFusioncharts
                type='mscolumn3d'
                width="100%"
                height="300%"
                dataFormat="JSON"
                dataSource={dataSource}
            />
        </div>
    );
    }
}
export default MyComponent2