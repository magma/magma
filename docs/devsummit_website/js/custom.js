(function($){
	"use strict";

	// Setup venobox / photo popup gallery
	$('.venobox').venobox({
		framewidth: '700px',
		frameheight: '400px',
		border: '10px',
		bgcolor: '#fff',
		titleattr: 'data-title',
		numeratio: true,
		infinigall: true
	});

	// Setup mobile menu
	$('.navbarmenumobile').on('click', function(){
		$(this).toggleClass('baropen');
		$('.navigation ul').toggleClass('open');
		return false;
	});


	// Setup masonry
	var $gridmain = $('.gridmain');
	$gridmain.masonry({
	    columnWidth: '.grid-sizer',
	    itemSelector: '.grid-item',
	    percentPosition: true
	});

	$gridmain.imagesLoaded().on('progress', function(){
	    $gridmain.masonry('layout');
	});


	// Setup banner slider
	var mySwiper = new Swiper ('.swiper-container', {
	    // Optional parameters
	    direction: 'horizontal',
	    loop: true,
	    autoplay: {
		    delay: 4000,
		},
	    navigation: {
	      nextEl: '.swiper-button-next',
	      prevEl: '.swiper-button-prev',
	    },
	});


    // Email Validation
	function isValidEmailAddress(emailAddress) {
		    var pattern = /^([a-z\d!#$%&'*+\-\/=?^_`{|}~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+(\.[a-z\d!#$%&'*+\-\/=?^_`{|}~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]+)*|"((([ \t]*\r\n)?[ \t]+)?([\x01-\x08\x0b\x0c\x0e-\x1f\x7f\x21\x23-\x5b\x5d-\x7e\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|\\[\x01-\x09\x0b\x0c\x0d-\x7f\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))*(([ \t]*\r\n)?[ \t]+)?")@(([a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|[a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF][a-z\d\-._~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]*[a-z\d\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])\.)+([a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]|[a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF][a-z\d\-._~\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]*[a-z\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])\.?$/i;
		    return pattern.test(emailAddress);
	};

    // Form Contact Validation
	$('form#ajax_form').on('submit', function(){
		var check = false;
		var name = $('#fname').val();
		var email = $('#femail').val();
		var company = $('#compname').val();

		//check name field
		if (name == ""){
			$('p.notif').text("Field name cannot be empty!").fadeIn();
			check = false;
			return false;
		}else{
			check = true;
		}

		//check email field
		if (email == ""){
			$('p.notif').text("Field email cannot be empty!").fadeIn();
			check = false;
			return false;
		}else{
			if( !isValidEmailAddress( email ) ) {
				$('p.notif').text("Email must be correct format!").fadeIn();
				check = false;
				return false;
			}else{
				check = true;
			}
		}

		//check phone field
		if (company == ""){
			$('p.notif').text("Field company cannot be empty!").fadeIn();
			check = false;
			return false;
		}else{
			check = true;
		}

		if (check == true){
			$("#btnsignup").prop('disabled', true);
			$("#btnsignup").prop('value', 'Sending in progress...');
			$.ajax({
				type: "POST",
				url: "postcontact.php",
				data: $('#ajax_form').serialize(),
				success: function(data){
					$('p.notif').html('<label>'+ data +'</label>').fadeIn();
					$("#btnsignup").prop('disabled', false);
					check = false;
					$('#fname').val("");
					$('#femail').val("");
					$('#compname').val("");
					$("#btnsignup").prop('value', 'Sign Up');
				}
			});
			return false;
		}
		return false;

	});


	// Smooth scroll
	$('a[href*="#"]')
	  .not('[href="#"]')
	  .not('[href="#0"]')
	  .click(function(event) {
	    if (
	      location.pathname.replace(/^\//, '') == this.pathname.replace(/^\//, '')
	      &&
	      location.hostname == this.hostname
	    ) {
	      var target = $(this.hash);
	      target = target.length ? target : $('[name=' + this.hash.slice(1) + ']');
	      if (target.length) {
	        event.preventDefault();
	        $('html, body').animate({
	          scrollTop: target.offset().top
	        }, 1000, function() {
	          var $target = $(target);
	          $target.focus();
	          if ($target.is(":focus")) {
	            return false;
	          } else {
	            $target.attr('tabindex','-1');
	            $target.focus();
	          };
	        });
	      }
	    }
	  });

})(jQuery);
