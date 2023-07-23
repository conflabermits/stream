package main

import (
	"fmt"
	"net/http"
)

/*
TO DO:
* Add function for checking current donorbox values
* Turn HTML content into a template
* Give the HTML template some inputs for changing reload interval, maybe other things
* Give the HTML template a button to save changes and reload with new values
*/

func main() {
	http.HandleFunc("/", serveHTML)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en" class="mdl-js"><head><style>.pac-container{background-color:#fff;position:absolute!important;z-index:1000;border-radius:2px;border-top:1px solid #d9d9d9;font-family:Arial,sans-serif;-webkit-box-shadow:0 2px 6px rgba(0,0,0,.3);box-shadow:0 2px 6px rgba(0,0,0,.3);-webkit-box-sizing:border-box;box-sizing:border-box;overflow:hidden}.pac-logo:after{content:"";padding:1px 1px 1px 0;height:18px;-webkit-box-sizing:border-box;box-sizing:border-box;text-align:right;display:block;background-image:url(https://maps.gstatic.com/mapfiles/api-3/images/powered-by-google-on-white3.png);background-position:right;background-repeat:no-repeat;-webkit-background-size:120px 14px;background-size:120px 14px}.hdpi.pac-logo:after{background-image:url(https://maps.gstatic.com/mapfiles/api-3/images/powered-by-google-on-white3_hdpi.png)}.pac-item{cursor:default;padding:0 4px;text-overflow:ellipsis;overflow:hidden;white-space:nowrap;line-height:30px;text-align:left;border-top:1px solid #e6e6e6;font-size:11px;color:#515151}.pac-item:hover{background-color:#fafafa}.pac-item-selected,.pac-item-selected:hover{background-color:#ebf2fe}.pac-matched{font-weight:700}.pac-item-query{font-size:13px;padding-right:3px;color:#000}.pac-icon{width:15px;height:20px;margin-right:7px;margin-top:6px;display:inline-block;vertical-align:top;background-image:url(https://maps.gstatic.com/mapfiles/api-3/images/autocomplete-icons.png);-webkit-background-size:34px 34px;background-size:34px}.hdpi .pac-icon{background-image:url(https://maps.gstatic.com/mapfiles/api-3/images/autocomplete-icons_hdpi.png)}.pac-icon-search{background-position:-1px -1px}.pac-item-selected .pac-icon-search{background-position:-18px -1px}.pac-icon-marker{background-position:-1px -161px}.pac-item-selected .pac-icon-marker{background-position:-18px -161px}.pac-placeholder{color:gray}.pac-target-input:-webkit-autofill{-webkit-animation-name:beginBrowserAutofill;animation-name:beginBrowserAutofill}.pac-target-input:not(:-webkit-autofill){-webkit-animation-name:endBrowserAutofill;animation-name:endBrowserAutofill}sentinel{}
</style><meta charset="utf-8"><link rel="stylesheet" href="https://donorbox.org/assets/application_donor-469a3c6c7121c1b24e02f3747a772a65b515b63996b27beceb36e09a36a5f682.css" media="all">

</head><body class="donor page--white" style="opacity: 1;">
  

  
<div class="main"><div class="hero with-org-admin-header" style="background-image: url(https://cdn.filestackcontent.com/QQfSVforQXCJEOtY3k3A)"><div class="overlay"></div><section><h1>Support Black Girls CODE</h1></section></div><div class="container content donor-content with-hero with-meter with-org-admin-header"><section class="new-page"><div class="row"><div class="large-7 columns description"><article id="fundraise_card" class="">

  <h3>
    Chris Dunaj
  </h3>

  <p>
    Fundraising on behalf of
    <a href="https://donorbox.org/support-black-girls-code">Black Girls CODE</a>
  </p>

  <picture id="avatar_picture" style="background-image: url(https://cdn.filestackcontent.com/fYNYJqS2TsqcUet2k7Uh)"></picture>
</article><style type="text/css">.progress .meter {
  width: 0%;
}</style><section class="donation-meter secondary-tabs"><div class="row collapse"><div class="small-12 columns"><div class="tabs-content" style="margin-bottom: 0"><div class="content active" id="panel-1"><div class="progress"><div class="meter"></div></div><div class="description"><div><p class="bold" id="total-raised">$0</p><p>Raised</p></div><div class="middle"><p class="bold" id="paid-count">0</p><p>Donations</p></div><div><p class="bold">$500</p><p>Goal</p></div></div></div></div></div></div></section><section class="content-description"><section class="donor-campaign-details fr-view" notranslate=""><p fr-original-style="" style="margin-top: 0px; margin-bottom: 0px;">Black Girls CODE builds pathways for Black girls to embrace the current tech marketplace as builders and creators by introducing them to skills in computer programming and technology, and closing the opportunity gap for Black women and girls. BGC is a global movement to establish equal representation in the tech sector and is devoted to showing the world that Black girls can code and do so much more. As someone who works in tech and sees the lower representation of both Black people and women in the field, I'm compelled to enable change that makes the next generation more diverse and equitable. Please support Black Girls CODE if you can!</p><p fr-original-style="" style="margin-top: 0px; margin-bottom: 0px;"><br fr-original-style="" style="margin-bottom: 0px;"></p><p fr-original-style="" style="margin-top: 0px; margin-bottom: 0px;"><span contenteditable="false" draggable="true" fr-original-class="fr-video fr-deletable fr-fvc fr-dvb fr-draggable" fr-original-style="" style="user-select: none; text-align: center; position: relative; display: block; clear: both;"><iframe width="640" height="360" src="https://www.youtube.com/embed/RSZ4kw0gusk?&amp;wmode=opaque" frameborder="0" allowfullscreen="" fr-original-style="" fr-original-class="fr-draggable video-blocker-visited" style="box-sizing: content-box; max-width: 100%; border: none;"></iframe></span><br fr-original-style="" style="margin-bottom: 0px;"></p></section><section class="desc-sharing mobile"><div id="sharing-buttons" class="sharing_buttons" notranslate="">
  <div class="ty-footer-content clearfix">
    <div class="dwc">
      <a role="button" class="resp-sharing-button__link" data-href="https://facebook.com/sharer/sharer.php?u=https%3A%2F%2Fdonorbox.org%2Fsupport-black-girls-code%2Ffundraiser%2Fchristopher-dunaj" aria-label="Facebook" target="_blank" onclick="off_x=(screen.width/2)-200;off_y=event.clientY-185;window.open(this.getAttribute('data-href'), 'targetWindow', 'toolbar=no,location=no,status=no,menubar=no,scrollbars=yes,resizable=yes,width=400,height=250,left='+off_x+',top='+off_y); return false;">
        <div class="resp-sharing-button resp-sharing-button--facebook resp-sharing-button--medium">
          <div aria-hidden="true" class="resp-sharing-button__icon resp-sharing-button__icon--solid">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path d="M18.77 7.46H14.5v-1.9c0-.9.6-1.1 1-1.1h3V.5h-4.33C10.24.5 9.5 3.44 9.5 5.32v2.15h-3v4h3v12h5v-12h3.85l.42-4z"></path>
            </svg>
          </div>
          Facebook
        </div>
      </a>
    </div>

    <div class="dwc">
      <a role="button" class="resp-sharing-button__link" data-href="https://www.linkedin.com/shareArticle?mini=true&amp;url=https%3A%2F%2Fdonorbox.org%2Fsupport-black-girls-code%2Ffundraiser%2Fchristopher-dunaj&amp;title=Support+Black+Girls+CODE" aria-label="LinkedIn" onclick="off_x=(screen.width/2)-300;off_y=event.clientY-360;window.open(this.getAttribute('data-href'), 'targetWindow', 'toolbar=no,location=no,status=no,menubar=no,scrollbars=yes,resizable=yes,width=600,height=600,left='+off_x+',top='+off_y); return false;">
        <div class="resp-sharing-button resp-sharing-button--linkedin resp-sharing-button--medium">
          <div aria-hidden="true" class="resp-sharing-button__icon resp-sharing-button__icon--solid" target="_blank">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path d="M6.5 21.5h-5v-13h5v13zM4 6.5C2.5 6.5 1.5 5.3 1.5 4s1-2.4 2.5-2.4c1.6 0 2.5 1 2.6 2.5 0 1.4-1 2.5-2.6 2.5zm11.5 6c-1 0-2 1-2 2v7h-5v-13h5V10s1.6-1.5 4-1.5c3 0 5 2.2 5 6.3v6.7h-5v-7c0-1-1-2-2-2z"></path>
            </svg>
          </div>
          LinkedIn
        </div>
      </a>
    </div>

    <div class="dwc">
      <a role="button" class="resp-sharing-button__link" data-href="https://twitter.com/intent/tweet/?text=Support+Black+Girls+CODE+%7C+Black+Girls+CODE&amp;url=https%3A%2F%2Fdonorbox.org%2Fsupport-black-girls-code%2Ffundraiser%2Fchristopher-dunaj" aria-label="Twitter" target="_blank" onclick="off_x=(screen.width/2)-300;off_y=event.clientY-285;window.open(this.getAttribute('data-href'), 'targetWindow', 'toolbar=no,location=no,status=no,menubar=no,scrollbars=yes,resizable=yes,width=600,height=450,left='+off_x+',top='+off_y); return false;">
        <div class="resp-sharing-button resp-sharing-button--twitter resp-sharing-button--medium">
          <div aria-hidden="true" class="resp-sharing-button__icon resp-sharing-button__icon--solid">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
              <path d="M23.44 4.83c-.8.37-1.5.38-2.22.02.93-.56.98-.96 1.32-2.02-.88.52-1.86.9-2.9 1.1-.82-.88-2-1.43-3.3-1.43-2.5 0-4.55 2.04-4.55 4.54 0 .36.03.7.1 1.04-3.77-.2-7.12-2-9.36-4.75-.4.67-.6 1.45-.6 2.3 0 1.56.8 2.95 2 3.77-.74-.03-1.44-.23-2.05-.57v.06c0 2.2 1.56 4.03 3.64 4.44-.67.2-1.37.2-2.06.08.58 1.8 2.26 3.12 4.25 3.16C5.78 18.1 3.37 18.74 1 18.46c2 1.3 4.4 2.04 6.97 2.04 8.35 0 12.92-6.92 12.92-12.93 0-.2 0-.4-.02-.6.9-.63 1.96-1.22 2.56-2.14z"></path>
            </svg>
          </div>
          Twitter
        </div>
      </a>
    </div>
  </div>
</div>
</section></section></div></div></section></div><div class="container comment-content"><section class="comments" id="donor-comments"></section></div></div><div id="fb-root"></div>  



  
  

  
<link rel="stylesheet" href="https://donorbox.org/assets/donation_page_pro-6aa5e5ec5e440fb380d793a6e7e2ea52174cd0ac431358e93cdcf58d599f8fbf.css" media="screen">  
<div class="custom-css-container"></div>
</body></html>
	`

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
}
