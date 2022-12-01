window.onload = getPosts();

async function getPosts(request = "/posts") {
  await fetch(request, {method:"post"})
    .then((response) => response.json())
    .then(function (json) {
      console.log(json)
      if(request.startsWith("/filtered?filter=postId")) {
        createPosts(json, true);
      } else {
        document.getElementById("container").innerHTML = "";
        createPosts(json);
      }
    });
  if (request != "/posts") {
    document.getElementById("load_more_btn").style.display = "none";
  } else {
    document.getElementById("load_more_btn").style.display = "";
  }
  addClickForPassword("signup_confirmpass", "signup_button");
  addClickForPassword("login_pass", "login_button");
}

function addClickForPassword(element, ele_id) {
  const ele = document.getElementById(element);
  if (ele) {
    ele.addEventListener("keyup", function(event) {
      if (event.key == "Enter") {
          document.getElementById(ele_id).click()
    };
  });
}
}

function reactionButton(postId, commentId, reactId, reactCount, setActive=false) {
  let returndata = "";
  var reactIcon = `fa-circle-down`;
  if(reactId == 1) {
    reactIcon = `fa-circle-up`;
  }

  let addClass = "";
  if(setActive) {
    addClass += "active"
  }

  var setStyle = ""
  if(commentId > 0) {
    setStyle += "height:60%;";
    addClass += " ps-1 py-0 pe-0";
  } else {
    addClass += " border";
  }

  returndata += `<button class="btn btn-sm btn-dark ${addClass}" style="${setStyle}" id="rbc${postId}_${commentId}_${reactId}" onclick="addReaction(${postId}, ${commentId}, ${reactId}, this)"><i class="fa-solid ${reactIcon}"></i>
                  <span class="badge text-info" id="rb${postId}_${commentId}_${reactId}">${reactCount}</span>
                  </button>`
  return returndata;
}

async function addReaction(postID, commentID, reactionID, targetButton) {
  await fetch(
    "/add_reaction?post_id=" +
      postID +
      "&comment_id=" +
      commentID +
      "&reaction_id=" +
      reactionID
  )
    .then((response) => response.json())
    .then(function (json) {
      console.log(json)
      for(let i=1; i < 3; i++) {
        document.getElementById("rb" + postID + "_" + commentID + "_" + i).innerHTML = json['rb'+i];
        document.getElementById("rbc" + postID + "_" + commentID + "_" + i).classList.remove("active")
      }
      if(json.userReaction > 0) {
        targetButton.classList.add("active");
      }
    });
}

async function signup() {
  const username = document.getElementById("signup_name");
  const email = document.getElementById("signup_email");
  const password = document.getElementById("signup_pass");
  const confirmPassword = document.getElementById("signup_confirmpass");
  let newUser = {
    name: username.value,
    email: email.value,
    password: password.value,
    confirmPassword: confirmPassword.value,
  };
  console.log(newUser);
  await fetch("/signup", {
    method: "POST",
    body: JSON.stringify(newUser),
  })
    .then((response) => response.json())
    .then((json) => {
      console.log(json);
      const modalHeading = document.getElementById("signup_result_heading");
      const modalBody = document.getElementById("signup_result_body");
      let modalBtn = document.getElementById("login");
      if (!json.status) {
        modalHeading.innerHTML = "Oh Snap!";
        modalBody.innerHTML = `${json.message}`;
        modalBtn.setAttribute("data-bs-target", "#signup_modal");
        modalBtn.textContent = "Sign up";
      } else {
        modalHeading.innerHTML = "Welcome!";
        modalBody.innerHTML = `You are now registered to Gritface!<br />
 You can now login.`;
        username.value = "";
        email.value = "";
        password.value = "";
        confirmPassword.value = "";
        modalBtn.setAttribute("data-bs-target", "#login_modal");
        modalBtn.textContent = "Log in";
      }
    });
}
async function login() {
  const email = document.getElementById("login_email");
  const password = document.getElementById("login_pass");
  let user = {
    email: email.value,
    password: password.value,
  };
  password.value = "";
  resetLoginModal();
  console.log(user);
  await fetch("/login", {
    method: "POST",
    body: JSON.stringify(user),
  })
    .then((response) => response.json())
    .then((json) => {
      if (!json.status) {
        const loginPassLabel = document.getElementById("login_pass_label");
        loginPassLabel.innerHTML = `password<br />${json.message}`;
      } else {
        console.log(`logged in successfully with uid ${json.message}`);
        const loginForm = document.getElementById("login_success");
        loginForm.submit();
      }
    });
}
async function resetLoginModal() {
  const loginPassLabel = document.getElementById("login_pass_label");
  const loginPass = document.getElementById("login_pass");
  loginPassLabel.innerHTML = "password";
  loginPass.value = "";
  await fetch("/checkSession")
    .then((response) => response.json())
    .then((json) => {
      if (json.status) {
        window.location.replace("/loginSuccess");
      } else {
        console.log("please sign up");
      }
    });
}

async function loadPosts() {
  if(document.getElementById("container").lastChild == undefined) {
    document.getElementById("load_more").textContent = "No more posts to load";
    return
  }

  const lastPostID = parseInt(document.getElementById("container").lastChild.id.substring(1))
  let lastPost = {
    lastPostID: lastPostID,
  };
  await fetch("/loadPosts", {
    method: "POST",
    body: JSON.stringify(lastPost),
  })
    .then((response) => response.json())
    .then((json) => {
      if(!json.length) {
        document.getElementById("load_more").textContent = "No more posts to load";
      } 
      createPosts(json)
    });
}

function createPostDiv(postUserPic, postUsername, postID, postHeading, postBody, postCats, postInsertTime, postUpdateTime, commentsLength, comments, likeNum, dislikeNum, userReaction) {
  let userRection1, userRection2 = false;
  if(userReaction == "1") {
    userRection1 =true;
  } else if(userReaction == "2") {
    userRection2 =true;
  }
  return `<section class="row" id="post_section">
  <div class="row">
    <div class="col-2 col-md-1 col-lg-1 mt-2">
      <img class="rounded-circle ms-2" style="border: 2px solid #54B4D3;" width="50" src="${postUserPic}">
    </div>
    <div class="col-7 mt-2">
      <h5 class="ps-xs-2 ms-sm-1 ps-sm-2 ps-md-4 ps-lg-2 text-start text-info">${postUsername}</h5>
    </div>
  </div>
<div data-bs-target="#collapse_post_comments${postID}" data-bs-toggle="collapse">
    <div class="text-white rounded my-2 py-2" id="post_div">
        <div class="col-11 offset-1 my-1" id="post_heading">
            <h4>${postHeading}</h4>
        </div>
        <div class="col-10 offset-1" id="post_body">
            <div class="border-top border-info bg-dark text-center" id="post_image"></div>
            <div class="text-justify my-2">
                <p>${postBody}</p>
            </div>
            <div class="text-secondary">
            <p>${postCats}</p>
            <div class="row text-secondary">
                <div class="col-6 order-0 text-left" id="post_insert_time">
                    <p>${postInsertTime}</p>
                </div>
                <div class="col-6 order-1 text-end" id="post_mod_time">
                    <p>${postUpdateTime}</p>
                </div>
            </div>
        </div>
    </div>
</div>
</div>

<!----- <div class="offset-lg-1 offset-md-1 offset-0 py-1"> ----->
  <!---  <div class="mx-4 mb-4 mb-lg-2 mb-md-2"> ---->
        <div class="row">
            <div class="col-10 offset-1" id="post_reactions_container${postID}">
                ${reactionButton(postID, 0, 1, likeNum, userRection1)}
                ${reactionButton(postID, 0, 2, dislikeNum, userRection2)}
              <div class=row>
                <p class="mx-1 pt-1 mb-2 text-info" id="number_of_comments"
                  data-bs-target="#collapse_post_comments${postID}" data-bs-toggle="collapse">
                  ${commentsLength}  <i class="fa-regular fa-comments pt-1" style="font-size:18px;"></i></p>
                  </div>
            </div>  
          </div>
        </div>
    <!---  </div> --->
      ${document.getElementById("user_name")? createCommentTextArea(postUserPic, postID):""}
  <div class="collapse" id="collapse_post_comments${postID}">
  ${comments}
  </div>
</section>`
}

function createCommentDiv(postID, commentID, commentUserPic, commentUsername, newComment, likeNumComment = 0, dislikeNumComment = 0, userReaction=0, update_time) {
  let userRection1, userRection2 = false;
  if(userReaction == "1") {
    userRection1 = true;
  } else if(userReaction == "2") {
    userRection2 = true;
  }
  return `
  <div class="row mx-auto pb-2" id="post_comments_container${postID}${commentID}">
    <div class="col-lg-10 mx-auto col-md-10 col-11 border rounded" style="background-color: #343a40;" id="post_comment_body${postID}${commentID}">
    
    <div class="row">
      <div class="col-2 col-lg-1 col-md-1 col-xl-1 pt-1">
        <img class="rounded-circle" style="border: 2px solid #54B4D3;" src="${commentUserPic}" width="50"><img>
        </div>
      <div class="col-10 col-lg-11 col-md-11 col-xl-11">
        <div class="row">
          <div class="col-12 col-md-6 col-lg-6 col-xl-6 ps-xs-4 ps-sm-0 ps-md-4 ps-lg-4 ps-lg-1 text-start">
            <h5 class="text-info pt-1 mb-0 pb-0 ps-xs-3 ps-0 ps-sm-3">${commentUsername}</h5>
            </div>
          <div class="col-12 col-md-6 col-lg-6 col-xl-6 text-start text-md-end">
            <p class="text-secondary" style="font-size: 0.8em;">${update_time}</p>
          </div>
        </div>
      
       <div class="col-12 word-wrap ps-xs-4 ps-sm-0 ps-md-4 ps-lg-4">
          <p class="mb-0 pb-0 text-light ps-md-2">${newComment}</p>
          </div>
      </div>
    </div>

  <div class="text-end pb-1 my-0" id="comment_reactions_container${postID}${commentID}">
  ${reactionButton(postID, commentID, 1, likeNumComment, userRection1)}
  ${reactionButton(postID, commentID, 2, dislikeNumComment, userRection2)}
  </div>
  </div>
  </div>`;
}

function createCommentTextArea(userPic, postID) {
  return `<div class="col-lg-10 col-md-10 col-11 mx-auto pb-1" id="user_comment">

  <!-- <div class="row">
  <div class="col-lg-2 col-md-2 d-none d-md-inline d-lg-inline">
  <img class="rounded-circle" style="max-width: 110%; border: 2px solid #54B4D3" src="${userPic}"></img>
  </div> -->

  
  <div class="input-group mb-2">
  <input type="text"
  class="form-control bg-dark border-info rounded-start text-light pt-1"
  id="newComment"
  style="resize:none; font-size: 0.8em;"
  placeholder="Write a comment">
  <button
  class="btn bg-info text-dark"
  type="button"
  style="width: 15%;"
  onclick="addComment(${postID})">
  <i class="fa-regular fa-comment"></i>
  </button>
  </div>
  </div>
  </div>
  </div>`
};

function createPosts(json, addToTop=false) {
  const container = document.getElementById("container");
  const userPic = document.getElementById("user_pic");
  for (const [key, postJSON] of Object.entries(json)) {
    // get categories
    let categories = "";
    if (postJSON.categories) {
      postJSON.categories.map(
        (category) => (categories += `#${category} `)
      );
    }
    // get like and dislike numbers
    let likeNum = 0,
      dislikeNum = 0;
    if (postJSON.reactions) {
      postJSON.reactions.map(function (reactions) {
        for (const [reaction_user_id, reaction] of Object.entries(
          reactions
        )) {
          if (reaction == "1") likeNum++;
          if (reaction == "2") dislikeNum++;
        }
      });
    }

    // create post div
    const postDiv = document.createElement("div");
    postDiv.classList.add("border", "rounded", "mx-auto", "col-lg-8", "col-md-10", "offset-sm-1", "col-sm-11", "col-12", "mt-2", "mb-4", "mb-lg-2", "mb-md-2");
    postDiv.id = "p" + postJSON.post_id;

    // loop and create divs of comments
    let comments = ``;
    let likeNumComment, dislikeNumComment;
    for (const [key, comment] of Object.entries(postJSON.comments)) {
      if (comment.profile_image == "")
        comment.profile_image = "static/images/raccoon_thumbnail7.jpg";
      likeNumComment = 0;
      dislikeNumComment = 0;
      if (comment.reactions) {
        comment.reactions.map(function (reactions) {
          for (const [key, reaction] of Object.entries(reactions)) {
            if (reaction == "1") likeNumComment++;
            if (reaction == "2") dislikeNumComment++;
          }
        });
      }
      comments += createCommentDiv(postJSON.post_id, comment.comment_id, comment.profile_image, comment.username, comment.body, likeNumComment, dislikeNumComment, comment.user_reaction, comment.insert_time.substring(0, 19));
    }
    
    // assemble the whole post div
    postDiv.innerHTML = createPostDiv(postJSON.profile_image, postJSON.username, postJSON.post_id, postJSON.heading, postJSON.body, categories, postJSON.insert_time, postJSON.update_time, Object.keys(postJSON.comments).length, comments, likeNum, dislikeNum, postJSON.user_reaction);
    
    if(addToTop) {
      container.prepend(postDiv);
    } else {
      container.append(postDiv);
    }
  }
}