import React from 'react';
import axios from 'axios';
import { Chart as ChartJS } from 'chart.js/auto';
import { Line } from 'react-chartjs-2';

class App extends React.Component {

  constructor(props) {
    console.log("================================================================================")
    console.log("App Constructor STARTS----------------------------------")
    super(props);
    this.state = {
      chartData : {
        labels: [...Array(40).keys()],
        datasets: [
          {
            id: 0,
            label: 'Label1',
            data: [69, 69, 69, 69],
          },
          {
            id: 1,
            label: 'Label2',
            data: [55, 55, 55, 55],
          },
        ],
        
      },
    };
    this.handleGraphChange = this.handleGraphChange.bind(this)
    console.log("App Constructor ENDS----------------------------------")
  }

  handleGraphChange(id, newChartData) {

    console.log("this.state.chartData: ", this.state.chartData)
    var output2 = {...this.state.chartData};
    var output = JSON.parse(JSON.stringify(this.state.chartData))
    console.log("output: ", output)
    console.log("output spread: ", output2)
    console.log("newChartData: ", newChartData);
    console.log("index operator: ", newChartData.datasets[0].data)
    
    setTimeout(() => { 
      if (id == 0) {
        output.datasets[0] = newChartData.datasets[0];
      } else if (id == 1) {
        output.datasets[1] = newChartData.datasets[0];
      }
      this.setState({
        chartData : output,
      });
      console.log("output after edit: ", output)
      console.log("this.state.chartData: ", this.state.chartData);
    }, 0);
  }

  render() {
    return (
      <div className="app">
        <div className='location1'>
          <div className='search'>
            <h2>Weather Comparison App</h2>
            <Location 
              id={0} 
              onGraphChange={this.handleGraphChange}
              initLocation='London'
            />
            <hr></hr>
            <Location 
              id={1}
              onGraphChange={this.handleGraphChange}
              initLocation='New York'
            />
            <hr></hr>
            <Line 
              data={this.state.chartData}
              height={250}
              width={400}
              options={{
                maintainAspectRatio: true, 
                responsive: false,
                scales: {
                  y: {title: {display: true, text: '°C',}},
                  x: {title: {display: true, text: 'Time',}}
                }
              }}
            />
          </div>
        </div>
      </div>
    );
  }
}

class Location extends React.Component {

  constructor(props) {
    console.log("Location Constructor STARTS----------------------------------");
    super(props);
    this.state = {
      id : this.props.id,
      locationName : '',
      locationTextField : '',
      country : '',
      temperatureC : '',
      temperatureF : '',
    };
    this.getCoords = this.getCoords.bind(this);
    this.getCoords(this.props.initLocation);
    console.log("Location Constructor ENDS----------------------------------");
  };

  getCoords(location) {
    console.log("getCoords STARTS----------------------------------");
    const geocodingURL = '/api/weather/' + location;
    
    axios.get(geocodingURL).then((response) => {
      console.log("AXIOS GET GEOCODINGURL STARTS---------------");
      console.log(response.data);
      console.log("AXIOS GET GEOCODINGURL ENDS---------------");

      // Updates city and temp
      this.setState({
        locationName : response.data.city_name,
        country : response.data.country_name,
        temperatureC : Math.round(response.data.city_temp),
        temperatureF : Math.round((response.data.city_temp) * 9 / 5 + 32),
      });

      // Updates graph
      var color = '';
      if (this.state.id == 0) {
        color = 'rgba(245, 179, 66, 1)';
      } else {
        color = 'rgba(66, 147, 245, 1)';
      }
      this.props.onGraphChange(
        this.state.id,
        {
          labels : [...Array(40).keys()],
          datasets : [{
            id : this.state.id,
            label : response.data.city_name + ', ' + response.data.country_name,
            data : response.data.forecast_temps,
            borderColor: color,
            
          }],
          
        }
      );
    });  

    console.log("getCoords ENDS----------------------------------");
  };

  handleForm = (event) => {
    event.preventDefault()
    this.setState({
       locationTextField : event.target.value
    })
    console.log("handleForm executed------")
  }

  handleAction = (event) => {
    event.preventDefault()
    this.getCoords(this.state.locationTextField)
    console.log("handleAction executed------")
  }

  render() {
    return(
      <div className='location'>
        <p>Enter the location:</p>
        <form onSubmit={this.handleAction}>
          <input type='text' value={this.state.locationTextField} 
          onChange={this.handleForm}/>
          <input type='submit'></input>
        </form>
        <p>{this.state.locationName}, {this.state.country}</p>
        <p>{this.state.temperatureC} °C</p>
        <p>{this.state.temperatureF} °F</p>

        

        {console.log("rendered")}
      </div>
    );
  };
}

export default App;
