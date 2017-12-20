$(function() {
  $('#shortener-form').on('submit', function(e) {
    if (!e.isDefaultPrevented()) {
      $('#message').html('');
      var url = "/api/register/";
      $.ajax({
        type: "POST",
        url: url,
        data: $(this).serialize(),
        success: function(data) {
          $('#message').html('Succesfully shortened to <a href="/api/r/' + data + '">' + data + '</a>.');
        }
      });
      return false;
    }
  })
});
