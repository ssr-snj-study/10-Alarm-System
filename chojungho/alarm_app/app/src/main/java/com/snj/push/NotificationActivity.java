package com.snj.push;


import android.app.Activity;
import android.app.NotificationManager;
import android.content.Context;
import android.os.Bundle;
import android.widget.TextView;


public class NotificationActivity extends Activity {
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_notification);
        CharSequence msg = "MainActivity에서 전달 받은 값 ";
        int id = 0;


        Bundle extras = getIntent().getExtras();
        if (extras == null) {
            msg = "error";
        } else {
            // 전달받은 인텐트의 id값을 확인해서 값을 가져옴
            id = extras.getInt("noti_Id");
        }

        // 전달 받은 값을 액티비티에 set
        TextView tv_ac_noti = (TextView) findViewById(R.id.tv_ac_noti);
        msg = msg + " " + id;
        tv_ac_noti.setText(msg);
        NotificationManager nm =
                (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);

        //노티피케이션 제거
        nm.cancel(id);
    }


}
