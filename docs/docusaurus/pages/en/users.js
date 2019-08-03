var NewComponent = React.createClass({
  render: function() {
    return (
      <div>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        {/* Meta  */}
        <meta name="description" content="Magma Dev Summit 2019" />
        <meta name="keywords" content="Magma, Facebook Connectivity" />
        <title>Magma Dev Summit 2019</title>
        <link href="css/bootstrap.min.css" rel="stylesheet" />
        <link href="css/fontawesome-all.min.css" rel="stylesheet" />
        <link href="css/venobox.css" rel="stylesheet" />
        <link href="css/swiper.min.css" rel="stylesheet" />
        <link href="css/style.css" rel="stylesheet" />
        {/* Banner */}
        <section id="herobanner" className="hero_section herobanner overlay">
          {/* Navigation menu */}
          <div className="container menu" id="signup">
            <div className="navigation logo" style={{background: 'url("css/logo.png") no-repeat center center', height: '80px'}}><h1><br /></h1></div>
            <div className="navigation nav">
              <div className="navbarmenumobile">
                <a href="#">
                  <div className="bar" />
                  <div className="bar" />
                </a>
              </div>
              <ul>
                <li><a href="#herobanner">Home</a></li>
                <li><a href="#what">What</a></li>
                <li><a href="#why">Why</a></li>
                <li><a href="#details">Details</a></li>
              </ul>
            </div>
          </div>
          <div className="container height-100">
            <div className="display-table">
              <div className="table-cell">
                <div className="col-md-7">
                  <div className="hero_content">
                    <h1 className="text-white">2019 Magma Dev Summit</h1>
                    <p>Come learn and hack with us.<br />Please RSVP by August 30th.</p>
                  </div>
                </div>
                <div className="col-md-5">
                  <div className="form_container">
                    <h2 className="text-white"><b>Interested in Attending</b></h2>
                    <div className="formbanner">
                      <form name="ajax_form" action id="ajax_form" method="POST" className="form-horizontal formreg mb-25">
                        <div className="form-group">
                          <div className="col-md-12">
                            <p><b> Date:</b> September 9th, 2019<br />
                              <b>Location:</b> Facebook HQ (Menlo Park, Ca)
                              <b>Meals:</b> Breakfast, Lunch, Happy Hour
                            </p>
                          </div>
                        </div>
                        <div className="btnsignup"><a href="https://docs.google.com/forms/d/e/1FAIpQLSe7leZ8CBFW_AHPVo6xtFqf49DGhB3_q8U_EvpWY-ZpuGk4HA/viewform?usp=sf_link"><div style={{color: 'white'}}>Sign Up</div></a></div>
                      </form>
                      <p className="notif" />
                      <p>Must be over 18 years of age</p>
                      <p>Additional details will be emailed to you</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
        {/* End Banner */}
        {/* About */}
        <section id="what" className="about_section about overlay">
          <div className="container">
            <div className="imgposabout" style={{background: 'url("css/dev_summit.jpg") no-repeat center center'}} />
            <div className="col-lg-4 col-lg-offset-4">
              <h2>What Is This Event</h2>
              <div className="line" />
            </div>
            <div className="col-lg-8" />
            <div className="col-lg-8 col-lg-offset-4">
              <p>This September, the Magma team from Facebook Connectivity (FBC) will host its Developer Summit at Facebook HQ (Menlo Park, California) for the first time after open sourcing the platform at Mobile World Congress\u{2019}19 in Barcelona earlier this year. It will be a day-long event for Magma developers, ecosystem partners and community members to showcase demos and presentations, in addition to participating in a Hackathon.  </p><br />
              <p>Throughout the day, you will encounter presentations that highlight different aspects of Magma; demonstrate Magma in action through many interesting use cases; and have an opportunity for hacking on the Magma platform in one of Facebook\u{2019}s labs.
              </p>
            </div>
          </div>
        </section>
        {/* End About */}
        {/* Why */}
        <section id="why" className="destination_section overlay">
          <div className="container">
            <div className="col-md-10 ">
              <h2>Why We're Hosting It</h2>
              <div className="line" />
              <div className="display-table wrapdestination">
                <div className="table-cell">
                  <div className="col-md-20">
                    <p>Magma Dev Summit creates an opportunity for developers interested in building solutions on Magma to meet the Magma team at Facebook and the platform\u{2019}s open source contributors. Throughout the summit, you will hear directly from Magma community members contributing to the project.</p><br />
                    <p>Facebook Connectivity is committed to bring more people online to a faster Internet. Magma is expected to bring more people online by enabling operators with open, flexible and extensible network solutions. At the heart of Magma lies a distributed Evolved Packet Core (EPC), a Federation Gateway, and a cloud-based Orchestrator. We are hosting this event to closely engage with the industry: system integrators, RAN vendors, EPC vendors, startups and other ecosystem partners.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
        {/* End Why */}
        {/* Why Cont */}
        <section id="whycont" className="activity_section">
          <div className="container">
            <div className="col-md-12">
              <h2>Why Attend</h2>
              <div className="line" />
              <p>Magma Dev Summit will be a great opportunity for mobile broadband and cellular network industry players to see Magma in action and learn how to contribute to and deploy Magma. Magma can be used for a wide variety of use cases across Fixed Wireless Access, Private LTE / CBRS, Carrier Wi-Fi, Mobile Broadband (4G / LTE, 5G), Massive and Industrial IoT, and the Network-as-a-Service.
              </p>
            </div>
          </div>
          <div className="display-table wrapdestination">
          </div>
        </section>
        {/* End Why Cont */}
        {/* Details */}
        <section id="details" className="benefit_section overlay">
          <div className="container">
            <div className="col-md-12 col-sm-12">
              <h2>Details</h2>
              <div className="line" />
            </div>
            <div className="col-md-3 col-sm-6">
              <div className="wrap">
                <img src="css/date.png" className="img-responsive overlay" alt="" />
                <h3>Event Date</h3>
                <p>September 9th,2019</p>
              </div>
            </div>
            <div className="col-md-3 col-sm-6">
              <div className="wrap">
                <img src="css/location.png" className="img-responsive overlay" alt="" />
                <h3>Location</h3>
                <p>Facebook Headquarters Menlo Park, Ca</p>
              </div>
            </div>
            <div className="col-md-3 col-sm-6">
              <div className="wrap">
                <img src="css/who.png" className="img-responsive overlay" alt="" />
                <h3>For Who</h3>
                <p>Must be over 18 years old </p>
              </div>
            </div>
            <div className="col-md-3 col-sm-6">
              <div className="wrap">
                <img src="css/time.png" className="img-responsive overlay" alt="" />
                <h3>Timing</h3>
                <p>All day event breakfast, lunch, happy hour</p>
              </div>
            </div>
          </div>
        </section>
        {/* End Details */}
        {/* Map */}
        <section id="mapview" className="mapview">
          <h3 />
          <div className="map" id="map" />
        </section>
        {/* End Map */}
        {/* End Sosmed */}
        {/* Footer */}
        <section id="footer" className="footer">
          <div className="container">
            <div className="col-md-4 col-sm-12">
              <div className="widget_content editContent contentmain">
                <h2 className="text-white">Magma</h2>
                <p>Learn about Magma</p>
                <p>at Facebook's HQ in Menlo Park Ca.</p>
              </div>
            </div>
            <div className="col-md-3 col-sm-4">
              <div className="widget_content editContent">
                <h3 className="text-white">Magma</h3>
                <ul>
                  <li><a href="#herobanner">Home</a></li>
                  <li><a href="#what">What</a></li>
                  <li><a href="#why">Why</a></li>
                  <li><a href="#details">Details</a></li>
                </ul>
              </div>
            </div>
            <div className="col-md-2 col-sm-4">
              <div className="widget_content editContent lastchild">
                <h3 className="text-white">Magma Links</h3>
                <ul>
                  <li><a href="https://connectivity.fb.com/">Facebook Connectivity</a></li>
                  <li><a href="https://facebookincubator.github.io/magma">GitHub</a></li>
                </ul>
              </div>
            </div>
            <div className="col-md-12 col-xs-12 cprights">
          <p className=\"text-center\">\u{00A9} 2019 Magma - https://facebookincubator.github.io/magma</p>
            </div>
          </div>
        </section>
        {/* End Footer */}
      </div>
    );
  }
});
