import React from 'react';
import ReactDOM from 'react-dom';
import MyComponent from './App.js'
import MyComponent2 from './timeline.js'
import './App.css';

class Home extends React.Component{
  constructor(props){
    super(props);
    this.state= {
      countries: [],
      countrycall: [],
      data: [],
      value: '',
      from: '',
      to: '',
      name: [],
      url: '/home',
      defurl: '/home',
      disp:'true'
    };
  }
  
  fetching(ur){
    fetch(ur)
		.then(response => response.json())
		.then(data => {
      var dataPoints = [];
      var x;
      var l=data.length;
      for(x in data[l-1]["Stats"]){
        if(x !== "TotalCases" && x!=="Date" && x!=="CurrentDayRecovered" && x !== "CurrentDayCases" && x !== "CurrentDayDeaths"){
          dataPoints.push({
              label: x,
              value: data[l-1]["Stats"][x]
          });
        }
      }
      this.setState({data: dataPoints, name: data[l-1]["Name"]});
		});
  }

  handleChange = (event) =>{
    this.setState({value : event.target.value, data: []});
    let ur;
    if(event.target.value){
      ur = '/country/' + event.target.value;
    }
    else{
      ur = '/home';
    }
      this.setState({url : ur});
      this.fetching(ur);
  }

  handleFromChange =(event) =>{
    this.setState({from:event.target.value})
  }

  handleToChange =(event) =>{
    this.setState({to:event.target.value})
  }  

  fetchdailydata(ur){
    fetch(ur)
		.then(response => response.json())
		.then(data => {
      this.setState({data: data, name: data[0]["Name"]});
		});
  }
  mySubmitHandler =(event) =>{
    console.log("hi in submit handler")
    this.setState({disp:'false', data: []})
    let ur;
    ur='/country/india?from=2020-05-01&to=2020-05-04'
    //ur='/country/'+this.state.value+'?from='+this.state.from+'&to='+this.state.to;
    this.setState({url : ur});
    this.fetchdailydata(ur);
  }

  componentDidMount(){
    var ar=[];
    Promise.all(
      [
        fetch('/countrynames').then(response => response.json()),
        fetch('/home').then(response => response.json())
      ]
    )
    .then(data => {
      for (var i = 0; i < data[0].length; i++) {
        ar.push(data[0][i]);
      }
      this.setState({countries : ar});
      let t =[];
      t.push(<option value=""></option>);
      for (i = 0; i < data[0].length; i++) {
        t.push(<option value={this.state.countries[i]}>{this.state.countries[i]}</option>) ;
      }
      this.setState({countrycall : t});
      var x;
      var dataPoints = [];
			var l=data[1].length;
      for(x in data[1][l-1]["Stats"]){
        if(x !== "TotalCases" && x!=="Date" && x!=="CurrentDayRecovered" && x !== "CurrentDayCases" && x !== "CurrentDayDeaths"){
          dataPoints.push({
              label: x,
              value: data[1][l-1]["Stats"][x]
          });
        }
      }
      this.setState({data: dataPoints, name: data[1][l-1]["Name"]});

    });
  }

  render(){
    return(
      <div>
       <h1>COVID-19</h1>
        <form>
          <p style={{textAlign: 'left'}}>Select Country     :
            <select value={this.state.value} onChange={this.handleChange}>
              {this.state.countrycall}
            </select>
          </p>      
        </form>
        <form style={{textAlign: 'left'}} onSubmit={this.mySubmitHandler}>
          <p><input type="text" placeholder="From" value={this.state.from} onChange={this.handleFromChange} /></p>
          <p><input type="text" placeholder="to" value={this.state.to} onChange={this.handleToChange}/></p>
          <input type="submit" />
        </form>
        { this.state.disp === 'true' ? <MyComponent data = {this.state.data} name = {this.state.name}/>:<MyComponent2 data = {this.state.data} name = {this.state.name} /> }
      </div>
    );
  }
}

ReactDOM.render(<Home />, document.getElementById('root'));