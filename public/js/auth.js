document.querySelectorAll('[data-password-toggle]').forEach(button=>{
  button.addEventListener('click',()=>{
    const input=document.getElementById(button.dataset.passwordToggle);
    if(!input)return;
    const visible=input.type==='text';
    input.type=visible?'password':'text';
    const icon=button.querySelector('i');
    if(icon)icon.className=visible?'fa-regular fa-eye':'fa-regular fa-eye-slash';
  });
});

const resetForm=document.getElementById('passwordResetForm');
const resetToken=document.getElementById('passwordResetToken');
const resetSubmit=document.getElementById('passwordResetSubmit');

if(resetForm&&resetToken){
  if(!resetToken.value){
    const params=new URLSearchParams(window.location.hash.replace(/^#/,''));
    const token=params.get('token')||'';
    if(/^[A-Za-z0-9_-]{40,128}$/.test(token)){
      resetToken.value=token;
      window.history.replaceState(null,'',window.location.pathname+window.location.search);
    }else if(resetSubmit){
      resetSubmit.disabled=true;
      resetSubmit.title='Abre el enlace de recuperación recibido por correo.';
    }
  }
}

document.querySelectorAll('[data-otp]').forEach(group=>{
  const inputs=[...group.querySelectorAll('input')];

  inputs.forEach((input,index)=>{
    input.addEventListener('input',()=>{
      input.value=input.value.replace(/\D/g,'').slice(-1);
      if(input.value && index<inputs.length-1)inputs[index+1].focus();
    });

    input.addEventListener('keydown',event=>{
      if(event.key==='Backspace' && !input.value && index>0)inputs[index-1].focus();
      if(event.key==='ArrowLeft' && index>0)inputs[index-1].focus();
      if(event.key==='ArrowRight' && index<inputs.length-1)inputs[index+1].focus();
    });

    input.addEventListener('paste',event=>{
      const digits=(event.clipboardData?.getData('text')||'').replace(/\D/g,'').slice(0,6);
      if(!digits)return;
      event.preventDefault();
      digits.split('').forEach((digit,i)=>{if(inputs[i])inputs[i].value=digit});
      inputs[Math.min(digits.length,inputs.length)-1]?.focus();
    });
  });
});
