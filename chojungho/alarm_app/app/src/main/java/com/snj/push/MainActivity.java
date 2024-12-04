package com.snj.push;

import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Build;
import android.os.Bundle;
import android.util.Log;
import android.widget.Button;
import android.widget.TextView;

import androidx.activity.EdgeToEdge;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.graphics.Insets;
import androidx.core.view.ViewCompat;
import androidx.core.view.WindowInsetsCompat;
import androidx.core.app.NotificationCompat;

import com.google.firebase.messaging.FirebaseMessaging;
import com.google.firebase.firestore.FirebaseFirestore;

//import android.app.Activity;
import android.app.Notification;
import android.app.NotificationManager;
import android.app.NotificationChannel;
import android.app.PendingIntent;
import android.content.Context;
import android.content.Intent;
import android.content.res.Resources;
import android.view.View;

import java.util.HashMap;
import java.util.Map;


public class MainActivity extends AppCompatActivity {
    private static final String TAG = "MainActivity";
    private TextView tokenTextView;
    private static final String CHANNEL_ID = MyFirebaseMessagingService.CHANNEL_ID;
    private static final String CHANNEL_NAME = MyFirebaseMessagingService.CHANNEL_NAME;
    private static final String CHANNEL_DESCRIPTION = MyFirebaseMessagingService.CHANNEL_DESCRIPTION;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        EdgeToEdge.enable(this);
        setContentView(R.layout.activity_main);
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main), (v, insets) -> {
            Insets systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars());
            v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom);
            return insets;
        });

        // UI 요소 초기화
        tokenTextView = findViewById(R.id.tokenTextView);
        Button getTokenButton = findViewById(R.id.getTokenButton);

        // 버튼 클릭 리스너
        getTokenButton.setOnClickListener(view -> fetchToken());

        btn_noti = (Button) findViewById(R.id.btn_noti);
        btn_noti.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View arg0) {
                NotificationActivity();
            }
        });
    }

    // FCM 토큰 가져오기
    private void fetchToken() {
        FirebaseMessaging.getInstance().getToken()
                .addOnCompleteListener(task -> {
                    if (!task.isSuccessful()) {
                        Log.w(TAG, "FCM 토큰 가져오기 실패", task.getException());
                        tokenTextView.setText("Failed to fetch token");
                        return;
                    }

                    // FCM 토큰
                    String token = task.getResult();
                    Log.d(TAG, "FCM 토큰: " + token);

                    // 토큰 출력
                    tokenTextView.setText(token);

                    // Firestore에 토큰 저장
                    saveTokenToFirestore(token);
                });
    }

    // Firestore에 토큰 저장
    private void saveTokenToFirestore(String token) {
        // Firestore 인스턴스 가져오기
        FirebaseFirestore db = FirebaseFirestore.getInstance();


        // Firestore에서 기존 토큰 확인
        db.collection("fcm_tokens")
                .whereEqualTo("device", token) // 동일 토큰 검색
                .get()
                .addOnCompleteListener(task -> {
                    if (task.isSuccessful() && !task.getResult().isEmpty()) {
                        // 토큰이 이미 Firestore에 존재함
                        Log.d(TAG, "토큰이 이미 존재합니다: " + token);
                    } else {
                        // 토큰이 Firestore에 없으므로 새로 추가
                        Map<String, Object> data = new HashMap<>();
                        String deviceName = android.os.Build.MODEL;
                        data.put("device", token);
                        data.put("user", deviceName);

                        // Firestore에 데이터 저장
                        db.collection("fcm_tokens")
                                .add(data) // Firestore에 추가
                                .addOnSuccessListener(documentReference ->
                                        Log.d(TAG, "토큰 저장 성공: " + documentReference.getId()))
                                .addOnFailureListener(e ->
                                        Log.w(TAG, "토큰 저장 실패", e));
                    }
                })
                .addOnFailureListener(e ->
                        Log.w(TAG, "토큰 확인 실패", e));
    }


    // 수신한 메시지 알림
    Button btn_noti = null;

    private void createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            int importance = NotificationManager.IMPORTANCE_HIGH;

            NotificationChannel channel = new NotificationChannel(CHANNEL_ID, CHANNEL_NAME, importance);
            channel.setDescription(CHANNEL_DESCRIPTION);

            NotificationManager notificationManager = getSystemService(NotificationManager.class);
            if (notificationManager != null) {
                notificationManager.createNotificationChannel(channel);
            }
        }
    }

    public void NotificationActivity() {
        createNotificationChannel();

        Resources res = getResources();

        Intent notificationIntent = new Intent(this, NotificationActivity.class);
        notificationIntent.putExtra("noti_Id", 9999); //전달할 값


        // PendingIntent: 위임된 인텐드 - NotificationBar에서 Bar를 클릭했을때 contentIntent를 전달한다.
        // contentIntent는 notificationIntent을 갖고 있으므로 9999가 전달 된다.
        PendingIntent contentIntent = PendingIntent.getActivity(
                this,
                0,
                notificationIntent,
                PendingIntent.FLAG_UPDATE_CURRENT | PendingIntent.FLAG_IMMUTABLE
        );

        NotificationCompat.Builder builder = new NotificationCompat.Builder(this, CHANNEL_ID);

        builder.setContentTitle("setContentTitle")
                .setContentText("setContentText")
                .setTicker("setTicker")
                .setSmallIcon(R.mipmap.ic_launcher)
                .setLargeIcon(BitmapFactory.decodeResource(res, R.mipmap.ic_launcher))
                .setContentIntent(contentIntent)             // 위에서 만든 PendingIntent 객체
                .setAutoCancel(true)                         // 터치시 노티가 지워짐
                .setWhen(System.currentTimeMillis())         // 누르면 바로 Noti를 실행
                .setDefaults(Notification.DEFAULT_ALL);      // SOUND, VIBRATE, LIGHTS 모두 실행

        // 빌드 버전이 롤리팝 이하일경우 실행.
        if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.LOLLIPOP) {
            builder.setCategory(Notification.CATEGORY_MESSAGE)
                    .setPriority(Notification.PRIORITY_HIGH);
//                    .setVisibility(Notification.VISIBILITY_PUBLIC);
        }

        NotificationManager nm = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);
        nm.notify(1234, builder.build());
    }

    // 수신한 메시지 알림
//    public void clickBtn(View view) {
//
//        //알림(Notification)을 관리하는 관리자 객체를 운영체제(Context)로부터 소환하기
//        NotificationManager notificationManager = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);
//
//        //Notification 객체를 생성해주는 건축가객체 생성(AlertDialog 와 비슷)
//        NotificationCompat.Builder builder = null;
//
//        //Oreo 버전(API26 버전)이상에서는 알림시에 NotificationChannel 이라는 개념이 필수 구성요소가 됨.
//        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
//
//            String channelID = "channel_01"; //알림채널 식별자
//            String channelName = "MyChannel01"; //알림채널의 이름(별명)
//
//            //알림채널 객체 만들기
//            NotificationChannel channel = new NotificationChannel(channelID, channelName, NotificationManager.IMPORTANCE_DEFAULT);
//
//            //알림매니저에게 채널 객체의 생성을 요청
//            notificationManager.createNotificationChannel(channel);
//
//            //알림건축가 객체 생성
//            builder = new NotificationCompat.Builder(this, channelID);
//
//
//        } else {
//            //알림 건축가 객체 생성
//            builder = new NotificationCompat.Builder(this, TAG);
//        }
//
//        //건축가에게 원하는 알림의 설정작업
//        builder.setSmallIcon(android.R.drawable.ic_menu_view);
//
//        //상태바를 드래그하여 아래로 내리면 보이는
//        //알림창(확장 상태바)의 설정
//        builder.setContentTitle("Title");//알림창 제목
//        builder.setContentText("Messages....");//알림창 내용
//        //알림창의 큰 이미지
////        Bitmap bm= BitmapFactory.decodeResource(getResources(),R.drawable.gametitle_09);
////        builder.setLargeIcon(bm);//매개변수가 Bitmap을 줘야한다.
//
//        //건축가에게 알림 객체 생성하도록
//        Notification notification = builder.build();
//
//        //알림매니저에게 알림(Notify) 요청
//        notificationManager.notify(1, notification);
//
//        //알림 요청시에 사용한 번호를 알림제거 할 수 있음.
//        //notificationManager.cancel(1);
//
//    }
}
