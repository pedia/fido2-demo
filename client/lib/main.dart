import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:flutter_web_auth/flutter_web_auth.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  final controller = TextEditingController(text: 'foo');

  String? _status;

  void onRegister() {
    final name = controller.text;
    final u = Uri.parse('https://9aedu.net/register');

    final q = Map<String, List<String>>.from(u.queryParametersAll);
    q['username'] = <String>[name];

    final u2 = Uri.https(u.authority, u.path, q);
    // launchUrl(u2);

    FlutterWebAuth.authenticate(url: u2.toString(), callbackUrlScheme: 'foobar')
        .then((result) {
      setState(() {
        _status = 'Got result: $result';
      });
    }).catchError((e) {
      setState(() {
        _status = 'Got error: $e';
      });
    });
  }

  void onLogin() {}

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            SelectableText(
              _status ?? '',
              style: Theme.of(context)
                  .textTheme
                  .bodyMedium
                  ?.apply(color: Colors.red),
            ),
            TextField(controller: controller),
            TextButton(
              onPressed: onRegister,
              child: Text(
                'Register',
                style: Theme.of(context).textTheme.headlineMedium,
              ),
            ),
            TextButton(
              onPressed: onLogin,
              child: Text(
                'Login',
                style: Theme.of(context).textTheme.headlineMedium,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
