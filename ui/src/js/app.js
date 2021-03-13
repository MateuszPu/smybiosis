import '../styles/index.scss';

if (process.env.NODE_ENV === 'development') {
  require('../index.html');
  console.log('start');
}

// Importing JavaScript
//
// You have two choices for including Bootstrap's JS files—the whole thing,
// or just the bits that you need.


// Option 1
//
// Import Bootstrap's bundle (all of Bootstrap's JS + Popper.js dependency)

// import "../../node_modules/bootstrap/dist/js/bootstrap.bundle.min.js";


// Option 2
//
// Import just what we need

// If you're importing tooltips or popovers, be sure to include our Popper.js dependency
// import "../../node_modules/popper.js/dist/popper.min.js";

import jQuery from "jquery/dist/jquery.slim";
import "bootstrap/js/dist/util";
import "bootstrap/js/dist/modal";
import "bootstrap/js/dist/collapse";

var mn = jQuery("#main-nav");
var shrinkNavbar = function() {
  mn.offset().top > 100 ? mn.addClass("navbar-shrink") : mn.removeClass("navbar-shrink");
};
jQuery(window).on('scroll', shrinkNavbar);
jQuery(window).on('load', function() {
  correctFaqHeight('#faq-checkboxes input:checked + label');
});

function correctFaqHeight(labelSelector) {
  // wait for toggle
  setTimeout(function() {
    var faqCheckboxes = jQuery('#faq-checkboxes');
    faqCheckboxes.height('auto');
    var answerHeight = jQuery(`${labelSelector} + p`).height();
    faqCheckboxes.height(answerHeight > faqCheckboxes.height() ? answerHeight : 'auto');
  }, 0);
}

window.scrollToFaq = function(labelSelector) {
  if (window.outerWidth > 991) { 
    // correct height because of height
    correctFaqHeight(labelSelector);
    return;
  }

  // wait for toggle
  setTimeout(function() {
    jQuery('html, body').scrollTop(jQuery(labelSelector).offset().top - 50);
  }, 0);
};

function isFormValid(form) {
  return form.checkValidity() === true && form.getElementsByClassName('.is-invalid').length === 0;
}

function toggleSubmitButton(form) {
  console.log(form, form.getElementsByTagName('button'), form.getElementsByClassName('.is-invalid').length);
  form.getElementsByTagName('button')[0].toggleAttribute('disabled', !isFormValid(form) );
}

/*
  REQUIRED STRUCTURE FOR VALIDATION:
  form .row .col .form-group input.form-control
*/

// JavaScript for disabling form submissions if there are invalid fields
(function() {
  'use strict';
  window.addEventListener('load', function() {
    var inputs = document.querySelectorAll('.needs-validation .form-group input:not(.custom-validation), .needs-validation .form-group select');
    var fieldValidation = Array.prototype.filter.call(inputs, function(input) {
      input.addEventListener('input', function(event) {
        input.parentNode.classList.add('was-validated');
        toggleSubmitButton(input.parentNode.parentNode.parentNode.parentNode);
      }, false);
    });

    var numberInputs = document.querySelectorAll('.needs-validation .form-group input.custom-validation.above-zero');
    var aboveZeroValidation = Array.prototype.filter.call(numberInputs, function(input) {
      input.addEventListener('input', function(event) {
        if (input.value <= 0) {
          input.classList.remove('is-valid');
          input.classList.add('is-invalid');
        } else {
          input.classList.add('is-valid');
          input.classList.remove('is-invalid');
        }
        input.parentNode.classList.add('was-validated');
        toggleSubmitButton(input.parentNode.parentNode.parentNode.parentNode);
      }, false);

      // allow only positive numbers
      input.onkeydown = function(event) {
        if (event.key == "-" || event.key == "e" || event.key == "E") {
          event.preventDefault();
          return false;
        }
      };
    });

    // Fetch all the forms we want to apply custom Bootstrap validation styles to
    var forms = document.getElementsByClassName('needs-validation');
    // Loop over them and prevent submission
    var validation = Array.prototype.filter.call(forms, function(form) {
      form.addEventListener('submit', function(event) {
        if (form.checkValidity() === false) {
          event.preventDefault();
          event.stopPropagation();
        }
        form.classList.add('was-validated');
      }, false);
    });
  }, false);
})();