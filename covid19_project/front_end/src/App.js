import React from 'react';
// import ReactDOM from 'react-dom';
import FusionCharts from "fusioncharts";
import charts from "fusioncharts/fusioncharts.charts";
import ReactFusioncharts from "react-fusioncharts";

// Resolves charts dependancy
charts(FusionCharts);

class MyComponent extends React.Component {
  constructor(props){
    super(props);
    this.state={
      type: 'pie2d',
    }
  }

  render() {
   var dataSource;

   if(this.state.type === 'pie2d'){
    dataSource = {
        chart: {
          caption: "Corona Data Chart Of",
          subcaption: this.props.name,
          plottooltext: "<b>$value</b> are $label",
          showlegend: "1",
          showpercentvalues: "0",
          legendposition: "bottom",
          usedataplotcolorforlabels: "1",
          theme: "fusion"
        },
        data: this.props.data,
      };
    }
    else if(this.state.type === 'pie3d'){
      dataSource = {
        chart: {
          caption: "Corona Data Chart Of",
          subcaption: this.props.name,
          showvalues: "1",
          showpercentintooltip: "0",
          numberprefix: "",
          enablemultislicing: "1",
          theme: "fusion"
        },
        data: this.props.data,
      };
    }
    else if(this.state.type === 'doughnut2d'){
      dataSource = {
        chart: {
          caption: "Corona Data Chart Of",
          subcaption: this.props.name,
          showpercentvalues: "0",
          defaultcenterlabel: "Coronavirus Cases",
          aligncaptionwithcanvas: "0",
          captionpadding: "0",
          decimals: "1",
          plottooltext:
            "<b>$percentValue</b> are <b>$label</b>",
          centerlabel: "$label: $value",
          theme: "fusion"
        },
        data: this.props.data,
      };
    }
    else if(this.state.type === 'doughnut3d'){
      dataSource = {
        chart: {
          caption: "Corona Data Chart Of",
          subcaption: this.props.name,
          enablesmartlabels: "1",
          showlabels: "1",
          numbersuffix: "",
          usedataplotcolorforlabels: "1",
          plottooltext: "$label, <b>$value</b>",
          theme: "fusion"
        },
        data: this.props.data,
      };
    }
    else if(this.state.type === 'line'){
      dataSource = {
        chart: {
          caption: "Corona Data Chart Of",
          yaxisname: "Coronavirus Cases",
          subcaption: this.props.name,
          numbersuffix: "",
          rotatelabels: "1",
          setadaptiveymin: "1",
          theme: "fusion"
        },
        data: this.props.data,
      };
    }
    else if(this.state.type === 'bar2d'){
      dataSource = {
        chart: {
          caption: "Corona Data Chart Of",
          subcaption: this.props.name,
          yaxisname: "Coronavirus Cases",
          aligncaptionwithcanvas: "0",
          plottooltext: "<b>$dataValue</b> cases",
          theme: "fusion"
        },
        data: this.props.data,
      };
    }
    else if(this.state.type === 'column2d'){
      dataSource = {
        chart: {
          caption: "Corona Data Chart Of",
          subcaption: this.props.name,
          yaxisname: "Coronavirus Cases",
          numbersuffix: "",
          theme: "fusion"
        },
        data: this.props.data,
      };
    }
    else if(this.state.type === 'column3d'){
      dataSource = {
        chart: {
          caption: "Corona Data Chart Of",
          subcaption: this.props.name,
          yaxisname: "Coronavirus{br}Cases",
          decimals: "1",
          theme: "fusion"
        },
        data: this.props.data,
      };
    }

    return (
      <div>
      <p style={{textAlign: 'left'}}>Select Chart Type    :
      <select value={this.state.type} onChange={(event)=>{this.setState({type: event.target.value})}}>
        <option value='pie2d'>Pie 2D</option>
        <option value='pie3d'>Pie 3D</option>
        <option value='doughnut2d'>Donut 2D</option>
        <option value='doughnut3d'>Donut 3D</option>
        <option value='line'>Simple Line</option>
        <option value='bar2d'>Simple Bar</option>
        <option value='column2d'>Simple Column</option>
        <option value='column3d'>Column 3D</option>
      </select>
      </p>
      <ReactFusioncharts
        type={this.state.type}
        width="100%"
        height="300%"
        dataFormat="JSON"
        dataSource={dataSource}
      />
      </div>
    );
  }
}

export default MyComponent